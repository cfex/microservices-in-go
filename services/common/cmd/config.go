package cmd

type AmqpConfig struct {
	PORT string
	HOST string
	ADDR string
	PWD  string
}

type config struct {
	amqp *AmqpConfig
}


