# Umarshal config style

In this example, the app config is a strongly typed structure.
Commands unmarshal the viper.Viper registery into that structure, 
which is then injected to app modules across the board.
