package woql

import "github.com/bdragon300/terminusgo/woql/schema"

func File(url string, options schema.FileOptions) schema.QueryResource {
	return getQueryResource(schema.Source{URL: url}, options)
}

func Remote(url string, options schema.FileOptions) schema.QueryResource {
	return getQueryResource(schema.Source{URL: url}, options)
}

func Post(url string, options schema.FileOptions) schema.QueryResource {
	return getQueryResource(schema.Source{Post: url}, options)
}

func getQueryResource(source schema.Source, options schema.FileOptions) schema.QueryResource {
	format := schema.FormatTypeCSV
	if options != nil {
		format = options.FileFormatType()
	}
	return schema.QueryResource{
		Source:  source,
		Format:  format,
		Options: options,
	}
}
