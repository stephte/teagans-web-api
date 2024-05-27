package migrations

import (
	"teagans-web-api/app/utilities/enctypes"
	"teagans-web-api/app/models"
	"gorm.io/gorm"
	"fmt"
)

func EncryptTaskData(db *gorm.DB) error {
	fmt.Println("Encrypting data!")
	var err error

	db.Transaction(func(tx *gorm.DB) error {
		var tasks []models.Task
		tx.Find(&tasks)
		fmt.Println(len(tasks))
		for _, task := range tasks {
			encTitle := enctypes.EncString(task.Title)
			encJson := enctypes.EncString(task.DetailJson)
			encHtml := enctypes.EncString(task.DetailHtml)
			
			fmt.Println("encTitle", encTitle)

			titleVal, err := encTitle.Value()
			if err != nil {
				return err
			}
			jsonVal, err := encJson.Value()
			if err != nil {
				return err
			}
			htmlVal, err := encHtml.Value()
			if err != nil {
				return err
			}

			fmt.Println("titleVal.(string)", titleVal.(string))

			task.Title = enctypes.EncString(titleVal.(string))
			task.DetailJson = enctypes.EncString(jsonVal.(string))
			task.DetailHtml = enctypes.EncString(htmlVal.(string))

			fmt.Println(task.Title)
			fmt.Println("")

			if err = tx.Save(&task).Error; err != nil {
				return err
			}
		}

		return nil
	})

	fmt.Println("Finished encrypting data!")

	return err
}