package command

func (c Command) Init(_ ...string) {
	err := c.ConfigApp.InitConfig()
	if err != nil {
		c.Logger.Error().Err(err).Msg("")
		return
	}
}
