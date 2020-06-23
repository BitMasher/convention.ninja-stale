package inventory

import (
	"convention.ninja/auth"
	"errors"
	"github.com/graphql-go/graphql"
)

var assetType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Asset",
	Description: "An asset",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
			Description: "The master id for this asset",
		},
		"category": &graphql.Field {
			Type: graphql.String,
			Description: "The category for this asset",
		},
		"model": &graphql.Field {
			Type: graphql.String,
			Description: "The model or type of asset",
		},
		"manufacturer": &graphql.Field {
			Type: graphql.String,
			Description: "The manufacturer of the asset",
		},
		"location": &graphql.Field {
			Type: graphql.String,
			Description: "The current known location of the asset",
		},
	},
})

func GetQuery(controller Controller) *graphql.Object {
	inventorySchema := graphql.NewObject(graphql.ObjectConfig{
		Name: "InventoryQueryApi",
		Description: "The inventory and asset management api",
		Fields: graphql.Fields{
			"assets": &graphql.Field {
				Name: "assets",
				Description: "Get the list of assets in the system",
				Type: graphql.NewList(assetType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					token := p.Context.Value("token")
					if token != nil && auth.ValidateToken("api", token.(string)) != nil {
						return controller.GetAssets(p.Context)
					}
					return nil, errors.New("invalid privileges")
				},
			},
			"search": &graphql.Field {
				Name: "search",
				Description: "Search assets by term or entities by barcode",
				Type:
			},
		},
	})

	return inventorySchema
}

func GetMutation(controller Controller) *graphql.Object {

}