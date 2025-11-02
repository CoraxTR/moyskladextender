package usecases

import (
	"mstorefgo/internal/order_processor"
	"mstorefgo/internal/xlsxbuilder"
)

func BuildUploadXlsx(builder *xlsxbuilder.XlsxBuilder, orders map[string]order_processor.ProcessedOrder) error {
	err := builder.Build(orders)
	if err != nil {
		return err
	}
	return nil
}
