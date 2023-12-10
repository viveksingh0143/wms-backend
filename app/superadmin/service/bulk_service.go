package service

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	adminService "star-wms/app/admin/service"
	"star-wms/app/base/dto/category"
	"star-wms/app/base/dto/customer"
	"star-wms/app/base/dto/machine"
	"star-wms/app/base/dto/product"
	baseService "star-wms/app/base/service"
	"star-wms/core/types"
	"star-wms/core/utils"
	"strconv"
	"strings"
)

type BulkService interface {
	ImportMaterialDataFromCSV(csvFilePath string) (bool, error)
	ImportProductDataFromCSV(csvFilePath string) (bool, error)
	ImportMachineDataFromCSV(plantCode string, csvFilePath string) (bool, error)
	ImportCustomerDataFromCSV(plantCode string, csvFilePath string) (bool, error)
	ImportPermissionDataFromCSV(csvFilePath string) (bool, error)
	ImportRoleDataFromCSV(csvFilePath string) (bool, error)
	ImportUserDataFromCSV(csvFilePath string) (bool, error)
	ImportCategoryDataFromCSV(csvFilePath string) (bool, error)
}

type DefaultBulkService struct {
	plantService      adminService.PlantService
	productService    baseService.ProductService
	categoryService   baseService.CategoryService
	machineService    baseService.MachineService
	customerService   baseService.CustomerService
	permissionService adminService.PermissionService
	roleService       adminService.RoleService
	userService       adminService.UserService
}

func NewBulkService(plantService adminService.PlantService, productService baseService.ProductService, categoryService baseService.CategoryService, machineService baseService.MachineService, customerService baseService.CustomerService, permissionService adminService.PermissionService, roleService adminService.RoleService, userService adminService.UserService) BulkService {
	return &DefaultBulkService{plantService: plantService, productService: productService, categoryService: categoryService, machineService: machineService, customerService: customerService, permissionService: permissionService, roleService: roleService, userService: userService}
}

func (s *DefaultBulkService) ImportMaterialDataFromCSV(csvFilePath string) (bool, error) {
	rows, err := utils.ProcessCsvFile(csvFilePath)
	if err != nil {
		return false, err
	}

	materialsCat, _, _, sfgCat, err := s.getOrCreateRootCategories()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get and create root categories")
		return false, err
	}

	uniqueCategories := make(map[string]bool)
	// Extract unique categories from rows
	for _, row := range rows {
		// Assume category information is in the 'material-type' column
		categoryName := row["material-type"]
		if categoryName != "" {
			uniqueCategories[categoryName] = true
		}
	}

	// Iterate over unique categories and update the database
	for categoryName := range uniqueCategories {
		categorySlug := utils.GenerateSlug(categoryName) // Replace with your slugify function
		if categorySlug == "sfg" || categorySlug == "fg" {
			continue
		}

		// Check and Create category if it doesn't exist
		exists := s.categoryService.ExistsBySlug(categorySlug, 0)
		// Create category if it doesn't exist
		if !exists {
			newCategory := &category.Form{
				Name:     categoryName,
				Slug:     categorySlug,
				Status:   types.StatusActive,
				Parent:   materialsCat,
				FullPath: fmt.Sprintf("%s%s/", materialsCat.FullPath, categorySlug),
			}
			err = s.categoryService.CreateCategory(newCategory)
			if err != nil {
				log.Error().Err(err).Msg("Failed to create category")
				return false, err
			}
		}
	}

	uniqueMaterials := make(map[string]*product.Form)
	// Extract unique categories from rows
	for _, row := range rows {
		materialCode := row["material-code"]
		materialDescription := row["material-description"]
		plant := row["plant"]
		uom := row["uom"]
		lineNumber := row["line-number"]
		componentItemNumber := row["component-item-number"]
		materialDescription2 := row["material-description2"]
		materialType := row["material-type"]
		componentQuantity := row["component-quantity"]
		unitOfMeasure := row["unit-of-measure"]

		if materialCode == "" || materialDescription == "" || plant == "" || uom == "" || lineNumber == "" || componentItemNumber == "" || materialDescription2 == "" || materialType == "" || componentQuantity == "" || unitOfMeasure == "" {
			continue
		}
		unitTypes := map[string]string{
			"PC": "PIECE",
			"GM": "WEIGHT",
			"KG": "WEIGHT",
			"MT": "WEIGHT",
			"LT": "LIQUID",
			"YD": "LENGTH",
			"SM": "LENGTH",
		}
		if unitTypes[uom] != "" {
			categorySlug := utils.GenerateSlug(materialType)
			var categoryForm *category.Form
			var productType string
			if categorySlug == "sfg" {
				categoryForm = sfgCat
				productType = "SEMI FINISHED GOODS"
			} else {
				categoryForm, _ = s.categoryService.GetCategoryBySlug(utils.GenerateSlug(materialType))
				productType = "RAW MATERIAL"
			}
			productSlug := utils.GenerateSlug(fmt.Sprintf("%s-%s", materialDescription2, componentItemNumber))
			uniqueMaterials[componentItemNumber] = &product.Form{
				ProductType:  productType,
				Name:         materialDescription2,
				Slug:         productSlug,
				Code:         componentItemNumber,
				CmsCode:      "",
				Description:  materialDescription2,
				UnitType:     unitTypes[unitOfMeasure],
				UnitWeight:   0,
				UnitValue:    unitOfMeasure,
				Status:       types.StatusActive,
				CategoryPath: categoryForm.FullPath,
				Category:     categoryForm,
			}
		}
	}

	// Iterate over unique categories and update the database
	for materialCode := range uniqueMaterials {
		uniqueMaterial := uniqueMaterials[materialCode]
		// Check and Create category if it doesn't exist
		exists := s.productService.ExistsByCode(uniqueMaterial.Code, 0)
		// Create category if it doesn't exist
		if !exists {
			err = s.productService.CreateProduct(uniqueMaterial)
			if err != nil {
				log.Error().Err(err).Msg("Failed to create material")
				return false, err
			}
		}
	}

	bomProducts := make(map[string]*product.Form)

	notFoundBomProduct := make(map[string]*product.Form)
	notFoundBomMaterial := make(map[string]*product.Form)
	// Extract unique categories from rows
	for _, row := range rows {
		plant := row["plant"]
		materialCode := row["material-code"]
		materialDescription := row["material-description"]
		lineNumber := row["line-number"]
		componentItemNumber := row["component-item-number"]
		materialDescription2 := row["material-description2"]
		componentQuantity := row["component-quantity"]

		if plant == "" || materialCode == "" || lineNumber == "" || componentItemNumber == "" || componentQuantity == "" {
			continue
		}
		if bomProducts[materialCode] == nil {
			productForm, err := s.productService.GetProductByCode(materialCode)
			if err != nil {
				notFoundBomProduct[materialCode] = &product.Form{
					Code: materialCode,
					Name: materialDescription,
				}
				log.Error().Err(err).Msgf("Failed to get product by code: %s", materialCode)
				continue
			} else {
				notFoundBomProduct[materialCode] = &product.Form{
					Code:        materialCode,
					Name:        materialDescription,
					ProductType: productForm.ProductType,
				}
			}
			bomProducts[materialCode] = productForm
		}

		ingredientForm, err := s.productService.GetProductByCode(componentItemNumber)
		if err != nil {
			notFoundBomMaterial[componentItemNumber] = &product.Form{
				Code: materialCode,
				Name: materialDescription2,
			}
			log.Error().Err(err).Msgf("Failed to get ingredient by code: %s", componentItemNumber)
			continue
		}
		if ingredientForm.ID <= 0 {
			log.Error().Err(err).Msgf("Failed to get ingredient by code: %s", componentItemNumber)
		}
		productForm := bomProducts[materialCode]
		ingredientQuantity, _ := strconv.ParseFloat(componentQuantity, 64)
		found := false
		for _, ingr := range productForm.Ingredients {
			if ingr.Ingredient.Code == ingredientForm.Code {
				ingr.Quantity = ingredientQuantity
				found = true
				break
			}
		}
		if !found {
			ingredientModel := &product.IngredientForm{
				Ingredient: ingredientForm,
				Quantity:   ingredientQuantity,
			}
			ingredientModel.IngredientID = ingredientForm.ID
			productForm.Ingredients = append(productForm.Ingredients, ingredientModel)
		}
	}

	for productCode := range bomProducts {
		bomProduct := bomProducts[productCode]
		err = s.productService.UpdateProduct(bomProduct.ID, bomProduct)
		if err != nil {
			log.Error().Err(err).Msg("Failed to update bom product")
			return false, err
		}
	}

	if len(notFoundBomProduct) > 0 {
		fmt.Println("Below Bom Products Not Found")
		headers := []string{"S.No.", "Code", "Name", "Type"}
		var data [][]string
		sNo := 1
		maxColumnWidths := make(map[int]int)

		// Find the maximum width for each column
		for i, header := range headers {
			maxColumnWidths[i] = len(header)
		}

		for _, itemInfo := range notFoundBomProduct {
			row := []string{fmt.Sprintf("%d", sNo), itemInfo.Code, itemInfo.Name, itemInfo.ProductType}
			data = append(data, row)
			for i, col := range row {
				if len(col) > maxColumnWidths[i] {
					maxColumnWidths[i] = len(col)
				}
			}
			sNo++
		}

		// Print underscores for headers
		for _, width := range maxColumnWidths {
			fmt.Printf("|%s", strings.Repeat("-", width+2)) // +2 for space padding
		}
		fmt.Println("|")

		// Print headers with dynamic width
		for i, h := range headers {
			fmt.Printf("| %-*s ", maxColumnWidths[i], h)
		}
		fmt.Println("|")

		// Print underscores for headers
		for _, width := range maxColumnWidths {
			fmt.Printf("|%s", strings.Repeat("-", width+2)) // +2 for space padding
		}
		fmt.Println("|")

		// Print data with dynamic width
		for _, row := range data {
			for i, col := range row {
				fmt.Printf("| %-*s ", maxColumnWidths[i], col)
			}
			fmt.Println("|")
		}

		// Print underscores for table end
		for _, width := range maxColumnWidths {
			fmt.Printf("|%s", strings.Repeat("-", width+2)) // +2 for space padding
		}
		fmt.Println("|")
	}

	if len(notFoundBomMaterial) > 0 {
		fmt.Println("Below Bom Materials Not Found")
		headers := []string{"S.No.", "Code", "Name"}
		var data [][]string
		sNo := 1
		for _, itemInfo := range notFoundBomMaterial {
			data = append(data, []string{fmt.Sprintf("%d", sNo), itemInfo.Code, itemInfo.Name})
			sNo++
		}
		// Print headers
		for _, h := range headers {
			fmt.Printf("%-10s", h) // Adjust the width (10 here) as needed
		}
		fmt.Println()
		// Print data
		for _, row := range data {
			for _, col := range row {
				fmt.Printf("%-10s", col) // Ensure width matches headers
			}
			fmt.Println()
		}
	}
	return true, nil
}

func (s *DefaultBulkService) ImportProductDataFromCSV(csvFilePath string) (bool, error) {
	rows, err := utils.ProcessCsvFile(csvFilePath)
	if err != nil {
		return false, err
	}
	_, _, fgCat, sfgCat, err := s.getOrCreateRootCategories()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get and create root categories")
		return false, err
	}

	uniqueProducts := make(map[string]*product.Form)
	for _, row := range rows {
		code := row["code"]
		description := row["description"]
		unitType := row["unittype"]
		unitWeight := row["unitweight"]
		unitWeightType := row["unitweighttype"]
		categoryType := row["category"]

		if code == "" || description == "" || unitType == "" || unitWeight == "" || unitWeightType == "" || categoryType == "" {
			continue
		}
		productType := ""
		var productCat *category.Form
		if categoryType == "FG" {
			productType = "FINISHED GOODS"
			productCat = fgCat
		} else if categoryType == "SFG" {
			productType = "SEMI FINISHED GOODS"
			productCat = sfgCat
		} else {
			continue
		}
		unitTypes := map[string]string{
			"PC": "PIECE",
			"GM": "WEIGHT",
			"KG": "WEIGHT",
			"MT": "WEIGHT",
			"LT": "LIQUID",
			"YD": "LENGTH",
			"SM": "LENGTH",
		}
		if unitType == "Pcs" {
			unitType = "PC"
		}
		if unitWeightType == "Grams" {
			unitWeightType = "GM"
		}
		if unitTypes[unitType] != "" {
			unitWeightValue, _ := strconv.ParseFloat(unitWeight, 64)
			productSlug := utils.GenerateSlug(fmt.Sprintf("%s-%s", description, code))
			uniqueProducts[code] = &product.Form{
				ProductType:  productType,
				Name:         description,
				Slug:         productSlug,
				Code:         code,
				CmsCode:      "",
				Description:  description,
				UnitType:     unitTypes[unitType],
				UnitWeight:   unitWeightValue,
				UnitValue:    unitWeightType,
				Status:       types.StatusActive,
				CategoryPath: productCat.FullPath,
				Category:     productCat,
			}
		}
	}

	// Iterate over unique categories and update the database
	for productCode := range uniqueProducts {
		uniqueProduct := uniqueProducts[productCode]
		// Check and Create category if it doesn't exist
		exists := s.productService.ExistsByCode(uniqueProduct.Code, 0)
		// Create category if it doesn't exist
		if !exists {
			err = s.productService.CreateProduct(uniqueProduct)
			if err != nil {
				log.Error().Err(err).Msg("Failed to create product")
				return false, err
			} else {
				log.Info().Msg("Product created successfully")
			}
		} else {
			log.Info().Msg("Product already exists")
		}
	}
	return true, nil
}

func (s *DefaultBulkService) getOrCreateRootCategories() (*category.Form, *category.Form, *category.Form, *category.Form, error) {
	productsCat := &category.Form{
		Name:     "Products",
		Slug:     "products",
		FullPath: "/products/",
		Status:   types.StatusActive,
	}
	if !s.categoryService.ExistsBySlug(productsCat.Slug, 0) {
		err := s.categoryService.CreateCategory(productsCat)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create root products category")
			return nil, nil, nil, nil, err
		}
	}

	productsCat, err := s.categoryService.GetCategoryBySlug(productsCat.Slug)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get products category")
		return nil, nil, nil, nil, err
	}

	fgCat := &category.Form{
		Name:   "Finished Goods",
		Slug:   "finished-goods",
		Status: types.StatusActive,
		Parent: productsCat,
	}
	fgCat.FullPath = fmt.Sprintf("%s%s/", productsCat.FullPath, fgCat.Slug)
	if !s.categoryService.ExistsBySlug(fgCat.Slug, 0) {
		err = s.categoryService.CreateCategory(fgCat)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create finished goods category")
			return nil, nil, nil, nil, err
		}
	}
	fgCat, err = s.categoryService.GetCategoryBySlug(fgCat.Slug)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get finished goods category")
		return nil, nil, nil, nil, err
	}

	sfgCat := &category.Form{
		Name:   "Semi Finished Goods",
		Slug:   "semi-finished-goods",
		Status: types.StatusActive,
		Parent: productsCat,
	}
	sfgCat.FullPath = fmt.Sprintf("%s%s/", productsCat.FullPath, sfgCat.Slug)

	if !s.categoryService.ExistsBySlug(sfgCat.Slug, 0) {
		err = s.categoryService.CreateCategory(sfgCat)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create semi finished goods category")
			return nil, nil, nil, nil, err
		}
	}
	sfgCat, err = s.categoryService.GetCategoryBySlug(sfgCat.Slug)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get semi finished goods category")
		return nil, nil, nil, nil, err
	}

	materialsCat := &category.Form{
		Name:     "Materials",
		Slug:     "materials",
		FullPath: "/materials/",
		Status:   types.StatusActive,
	}
	if !s.categoryService.ExistsBySlug(materialsCat.Slug, 0) {
		err := s.categoryService.CreateCategory(materialsCat)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create raw material category")
			return nil, nil, nil, nil, err
		}
	}

	materialsCat, err = s.categoryService.GetCategoryBySlug(materialsCat.Slug)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get raw material category")
		return nil, nil, nil, nil, err
	}
	return materialsCat, productsCat, fgCat, sfgCat, nil
}

func (s *DefaultBulkService) ImportMachineDataFromCSV(plantCode string, csvFilePath string) (bool, error) {
	if plantCode == "" {
		log.Fatal().Msg("Plant code needed to import")
		return false, errors.New("plant code needed to import")
	}
	plant, err := s.plantService.GetPlantByCode(plantCode)
	if err != nil {
		return false, err
	}

	rows, err := utils.ProcessCsvFile(csvFilePath)
	if err != nil {
		return false, err
	}
	uniqueMachines := make(map[string]*machine.Form)
	for _, row := range rows {
		name := row["name"]
		code := row["code"]
		status := types.StatusActive
		switch strings.ToLower(row["status"]) {
		case "inactive":
			status = types.StatusInactive
		case "banned":
			status = types.StatusBanned
		}
		if name == "" || code == "" {
			continue
		}
		uniqueMachines[code] = &machine.Form{
			Name:   name,
			Code:   code,
			Status: status,
		}
	}

	// Iterate over unique categories and update the database
	for machineCode := range uniqueMachines {
		uniqueMachine := uniqueMachines[machineCode]
		//exists := s.machineService.ExistsByCode(plant.ID, uniqueMachine.Code, 0)
		machineData, err := s.machineService.GetMachineByCode(plant.ID, uniqueMachine.Code)
		if err != nil {
			err = s.machineService.CreateMachine(plant.ID, uniqueMachine)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to create machine, whose code is: %s", uniqueMachine.Code)
			} else {
				log.Info().Msg("Machine created successfully")
			}
		} else {
			err = s.machineService.UpdateMachine(plant.ID, machineData.ID, uniqueMachine)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to update machine, whose code is: %s", uniqueMachine.Code)
			} else {
				log.Info().Msg("Machine updated successfully")
			}
		}
		//if !exists {
		//	err = s.machineService.CreateMachine(plant.ID, uniqueMachine)
		//	if err != nil {
		//		log.Error().Err(err).Msg("Failed to create machine")
		//		return false, err
		//	} else {
		//		log.Info().Msg("Machine created successfully")
		//	}
		//} else {
		//	log.Info().Msg("Machine already exists")
		//}
	}
	return true, nil
}

func (s *DefaultBulkService) ImportCustomerDataFromCSV(plantCode string, csvFilePath string) (bool, error) {
	if plantCode == "" {
		log.Fatal().Msg("Plant code needed to import")
		return false, errors.New("plant code needed to import")
	}
	plant, err := s.plantService.GetPlantByCode(plantCode)
	if err != nil {
		return false, err
	}

	rows, err := utils.ProcessCsvFile(csvFilePath)
	if err != nil {
		return false, err
	}
	uniqueCustomers := make(map[string]*customer.Form)
	for _, row := range rows {
		name := row["name"]
		code := row["code"]
		contactperson := row["contactperson"]
		billingaddress1 := row["billingaddress1"]
		billingaddress2 := row["billingaddress2"]
		billingstate := row["billingstate"]
		billingcountry := row["billingcountry"]
		billingpincode := row["billingpincode"]
		shippingaddress1 := row["shippingaddress1"]
		shippingaddress2 := row["shippingaddress2"]
		shippingstate := row["shippingstate"]
		shippingcountry := row["shippingcountry"]
		shippingpincode := row["shippingpincode"]
		status := types.StatusActive
		switch strings.ToLower(row["status"]) {
		case "inactive":
			status = types.StatusInactive
		case "banned":
			status = types.StatusBanned
		}

		if name == "" || code == "" {
			continue
		}

		uniqueCustomers[code] = &customer.Form{
			Name:             name,
			Code:             code,
			ContactPerson:    contactperson,
			BillingAddress1:  billingaddress1,
			BillingAddress2:  billingaddress2,
			BillingState:     billingstate,
			BillingCountry:   billingcountry,
			BillingPincode:   billingpincode,
			ShippingAddress1: shippingaddress1,
			ShippingAddress2: shippingaddress2,
			ShippingState:    shippingstate,
			ShippingCountry:  shippingcountry,
			ShippingPincode:  shippingpincode,
			Status:           status,
		}
	}

	// Iterate over unique categories and update the database
	for customerCode := range uniqueCustomers {
		uniqueCustomer := uniqueCustomers[customerCode]
		//exists := s.customerService.ExistsByCode(plant.ID, uniqueCustomer.Code, 0)
		customerData, err := s.customerService.GetCustomerByCode(plant.ID, uniqueCustomer.Code)
		if err != nil {
			err = s.customerService.CreateCustomer(plant.ID, uniqueCustomer)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to create customer, whose code is: %s", uniqueCustomer.Code)
			} else {
				log.Info().Msg("Customer created successfully")
			}
		} else {
			err = s.customerService.UpdateCustomer(plant.ID, customerData.ID, uniqueCustomer)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to update customer, whose code is: %s", uniqueCustomer.Code)
			} else {
				log.Info().Msg("Customer updated successfully")
			}
		}
	}
	return true, nil
}

func (s *DefaultBulkService) ImportPermissionDataFromCSV(csvFilePath string) (bool, error) {
	return true, nil
}

func (s *DefaultBulkService) ImportRoleDataFromCSV(csvFilePath string) (bool, error) {
	return true, nil
}

func (s *DefaultBulkService) ImportUserDataFromCSV(csvFilePath string) (bool, error) {
	return true, nil
}

func (s *DefaultBulkService) ImportCategoryDataFromCSV(csvFilePath string) (bool, error) {
	return true, nil
}
