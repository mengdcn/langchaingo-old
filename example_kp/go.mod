module example_kp

go 1.20

require github.com/tmc/langchaingo v0.0.0-20230829032728-c85d3967da08

require (
	github.com/dlclark/regexp2 v1.8.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/pkoukk/tiktoken-go v0.1.2 // indirect
)

//replace github.com/tmc/langchaingo v0.0.0-20230829032728-c85d3967da08 => github.com/comqositi/langchaingo v0.0.0-20230830073207-b3bd1db2b0f2
replace github.com/tmc/langchaingo v0.0.0-20230829032728-c85d3967da08 => /home/lucas/bgy/langchaingo
