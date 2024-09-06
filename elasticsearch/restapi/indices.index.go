package restapi

func (i *Indices) Settings(index string) (string, error) {
	resp, err := i.client.Indices.GetSettings(
		i.client.Indices.GetSettings.WithIndex(index),
		i.client.Indices.GetSettings.WithIncludeDefaults(true),
	)
	if err != nil {
		return "", err
	}

	return GetResponseBody(resp, err)
}

func (i *Indices) Mapping(index string) (string, error) {
	resp, err := i.client.Indices.GetMapping(i.client.Indices.GetMapping.WithIndex(index))
	if err != nil {
		return "", err
	}

	return GetResponseBody(resp, err)
}

func (i *Indices) Stats(index string) {
	_, err := i.client.Indices.Stats(i.client.Indices.Stats.WithIndex(index))
	if err != nil {

	}
}
