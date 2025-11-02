package xlsxbuilder

import (
	"fmt"
	"mstorefgo/internal/config"
	"mstorefgo/internal/order_processor"
	"strconv"
	"time"

	"codeberg.org/tealeg/xlsx/v4"
)

type XlsxBuilder struct {
	cfg config.Config
}

func NewXlsxBuilder(cfg config.Config) *XlsxBuilder {
	return &XlsxBuilder{
		cfg: cfg,
	}
}
func clickCell(sh *xlsx.Sheet, row, col int) *xlsx.Cell {
	theCell, err := sh.Cell(row, col)
	if err != nil {
		panic(err)
	}
	return theCell
}

func BuildUploadXlsx(cfg config.Config, orders map[string]order_processor.ProcessedOrder) error {

	uploadFile, err := xlsx.OpenFile("../blankimport.xlsx")
	if err != nil {
		panic(err)
	}
	var importSheet *xlsx.Sheet
	var summarySheet *xlsx.Sheet
	var theCell *xlsx.Cell

	importSheet = uploadFile.Sheet["Импорт"]
	summarySheet = uploadFile.Sheet["Расписка"]

	importrownumber := 1
	summaryrownumber := 6

	for _, order := range orders {
		refnumber, err := strconv.Atoi(order.RefGoNumber)
		if err != nil {
			fmt.Println(err)
		}
		theCell = clickCell(importSheet, importrownumber, 0)
		theCell.SetInt(refnumber)
		theCell = clickCell(importSheet, importrownumber, 1)
		theCell.SetInt(refnumber)
		theCell = clickCell(importSheet, importrownumber, 2)
		theCell.SetString(order.RecieverName)
		theCell = clickCell(importSheet, importrownumber, 3)
		theCell.SetInt(order.RecieverPhoneNumber)
		theCell = clickCell(importSheet, importrownumber, 4)
		theCell.SetString(order.ShipmentAddress)
		theCell = clickCell(importSheet, importrownumber, 5)
		theCell.SetString(order.DeliveryPlannedDate)
		theCell = clickCell(importSheet, importrownumber, 6)
		theCell.SetString(order.DeliveryIntervalFrom)
		theCell = clickCell(importSheet, importrownumber, 7)
		theCell.SetString(order.DeliveryIntervalUntil)
		theCell = clickCell(importSheet, importrownumber, 8)
		theCell.SetString(order.Description)

		if order.ChilledBoxes > 0 && order.FrozenBoxes > 0 {
			theCell = clickCell(importSheet, importrownumber, 9)
			theCell.SetString("Охлаждённая продукция")
			theCell = clickCell(importSheet, importrownumber, 10)
			theCell.SetInt(refnumber)
			theCell = clickCell(importSheet, importrownumber, 11)
			theCell.SetInt(int(order.ChilledBoxes))

			if order.PaymentMethod == "Наличные" || order.PaymentMethod == "Терминал" {
				theCell = clickCell(importSheet, importrownumber, 12)
				theCell.SetFloat(order.Sum)
				theCell = clickCell(importSheet, importrownumber, 13)
				theCell.SetFloat(order.Sum)
			} else {
				theCell = clickCell(importSheet, importrownumber, 12)
				theCell.SetFloat(0)
				theCell = clickCell(importSheet, importrownumber, 13)
				theCell.SetFloat(0)
			}

			theCell = clickCell(importSheet, importrownumber, 14)
			theCell.SetFloat(order.ChilledWeight)
			theCell = clickCell(importSheet, importrownumber, 15)
			theCell.SetString("Средние температуры (+2+6)")

			theCell = clickCell(importSheet, importrownumber, 16)
			if order.PaymentMethod == "расч. счет" {
				theCell.SetString("Да")
			} else {
				theCell.SetString("Нет")
			}

			theCell = clickCell(importSheet, importrownumber, 17)
			theCell.SetString(order.DeliveryRegion)
			theCell = clickCell(importSheet, importrownumber, 18)
			theCell.SetInt(int(order.ChilledBoxes))
			importrownumber++
			theCell = clickCell(importSheet, importrownumber, 9)
			theCell.SetString("Замороженная продукция")
			theCell = clickCell(importSheet, importrownumber, 10)
			theCell.SetInt(refnumber)
			theCell = clickCell(importSheet, importrownumber, 11)
			theCell.SetInt(int(order.FrozenBoxes))
			theCell = clickCell(importSheet, importrownumber, 12)
			theCell.SetFloat(0)
			theCell = clickCell(importSheet, importrownumber, 13)
			theCell.SetFloat(0)
			theCell = clickCell(importSheet, importrownumber, 14)
			theCell.SetFloat(order.FrozenWeight)
			theCell = clickCell(importSheet, importrownumber, 15)
			theCell.SetString("Низкие температуры (-18)")
			theCell = clickCell(importSheet, importrownumber, 16)
			if order.PaymentMethod == "расч. счет" {
				theCell.SetString("Да")
			} else {
				theCell.SetString("Нет")
			}
			theCell = clickCell(importSheet, importrownumber, 17)
			theCell.SetString(order.DeliveryRegion)
			theCell = clickCell(importSheet, importrownumber, 18)
			theCell.SetInt(int(order.FrozenBoxes))
		} else {
			if order.ChilledBoxes == 0 {
				theCell = clickCell(importSheet, importrownumber, 9)
				theCell.SetString("Замороженная продукция")
				theCell = clickCell(importSheet, importrownumber, 10)
				theCell.SetInt(refnumber)
				theCell = clickCell(importSheet, importrownumber, 11)
				theCell.SetInt(int(order.FrozenBoxes))

				if order.PaymentMethod == "Наличные" || order.PaymentMethod == "Терминал" {
					theCell = clickCell(importSheet, importrownumber, 12)
					theCell.SetFloat(order.Sum)
					theCell = clickCell(importSheet, importrownumber, 13)
					theCell.SetFloat(order.Sum)
				} else {
					theCell = clickCell(importSheet, importrownumber, 12)
					theCell.SetFloat(0)
					theCell = clickCell(importSheet, importrownumber, 13)
					theCell.SetFloat(0)
				}

				theCell = clickCell(importSheet, importrownumber, 14)
				theCell.SetFloat(order.FrozenWeight)
				theCell = clickCell(importSheet, importrownumber, 15)
				theCell.SetString("Низкие температуры (-18)")

				theCell = clickCell(importSheet, importrownumber, 16)
				if order.PaymentMethod == "расч. счет" {
					theCell.SetString("Да")
				} else {
					theCell.SetString("Нет")
				}

				theCell = clickCell(importSheet, importrownumber, 17)
				theCell.SetString(order.DeliveryRegion)
				theCell = clickCell(importSheet, importrownumber, 18)
				theCell.SetInt(int(order.FrozenBoxes))
			} else {
				theCell = clickCell(importSheet, importrownumber, 9)
				theCell.SetString("Охлаждённая продукция")
				theCell = clickCell(importSheet, importrownumber, 10)
				theCell.SetInt(refnumber)
				theCell = clickCell(importSheet, importrownumber, 11)
				theCell.SetInt(int(order.ChilledBoxes))

				if order.PaymentMethod == "Наличные" || order.PaymentMethod == "Терминал" {
					theCell = clickCell(importSheet, importrownumber, 12)
					theCell.SetFloat(order.Sum)
					theCell = clickCell(importSheet, importrownumber, 13)
					theCell.SetFloat(order.Sum)
				} else {
					theCell = clickCell(importSheet, importrownumber, 12)
					theCell.SetFloat(0)
					theCell = clickCell(importSheet, importrownumber, 13)
					theCell.SetFloat(0)
				}

				theCell = clickCell(importSheet, importrownumber, 14)
				theCell.SetFloat(order.ChilledWeight)
				theCell = clickCell(importSheet, importrownumber, 15)
				theCell.SetString("Средние температуры (+2+6)")

				theCell = clickCell(importSheet, importrownumber, 16)
				if order.PaymentMethod == "расч. счет" {
					theCell.SetString("Да")
				} else {
					theCell.SetString("Нет")
				}

				theCell = clickCell(importSheet, importrownumber, 17)
				theCell.SetString(order.DeliveryRegion)
				theCell = clickCell(importSheet, importrownumber, 18)
				theCell.SetInt(int(order.ChilledBoxes))
			}
		}

		theCell = clickCell(summarySheet, summaryrownumber, 11)
		theCell.SetInt(refnumber)
		theCell = clickCell(summarySheet, summaryrownumber, 12)
		theCell.SetInt(int(order.ChilledBoxes) + int(order.FrozenBoxes))
		theCell = clickCell(summarySheet, summaryrownumber, 13)
		theCell.SetFloat(order.Sum)

		importrownumber++
		summaryrownumber++
	}

	temptoday := time.Now()
	today := temptoday.Format("02.01.2006")
	savepath := "../" + today + ".xlsx"

	uploadFile.Save(savepath)
	importSheet.Close()
	return nil
}

func (x *XlsxBuilder) Build(orders map[string]order_processor.ProcessedOrder) error {
	err := BuildUploadXlsx(x.cfg, orders)
	if err != nil {
		return err
	}
	return nil
}
