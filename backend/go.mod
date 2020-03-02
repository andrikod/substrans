module subtrans.com/subtrans

go 1.13

require (
	cloud.google.com/go v0.53.0
	golang.org/x/text v0.3.2
	subtrans.com/subtrans/handler v0.0.0
)

// Import locally the function
replace subtrans.com/subtrans/handler => ./functions/subTranslate
