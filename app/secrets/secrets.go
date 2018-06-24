package secrets

//go:generate safekeeper --output=appsecrets.go --keys=LOCAL_CLIENT_ID,LOCAL_CLIENT_SECRET,PROD_CLIENT_ID,PROD_CLIENT_SECRET,TEST_STRIPE_KEY,TEST_STRIPE_PUBLISHABLE_KEY,PROD_STRIPE_KEY,PROD_STRIPE_PUBLISHABLE_KEY $GOFILE
