package firebase

// [START admin_import_golang]
import (
	"context"
	"log"

	//firebase "firebase.google.com/go"
	firebase "firebase.google.com/go"
	//"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// [END admin_import_golang]

// ==================================================================
// https://firebase.google.com/docs/admin/setup
// ==================================================================

//InitializeAppWithServiceAccount init fcm
func InitializeAppWithServiceAccount() *firebase.App {
	// [START initialize_app_service_account_golang]
	opt := option.WithCredentialsFile("./runex-f341b-firebase-adminsdk-o6pr0-0757998bf7.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println("error initializing app: %v\n", err)
	}
	// [END initialize_app_service_account_golang]

	return app
}

func initializeAppWithRefreshToken() *firebase.App {
	// [START initialize_app_refresh_token_golang]
	opt := option.WithCredentialsFile("./runex-f341b-firebase-adminsdk-o6pr0-0757998bf7.json")
	config := &firebase.Config{ProjectID: "my-project-id"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Println("error initializing app: %v\n", err)
	}
	// [END initialize_app_refresh_token_golang]

	return app
}

func initializeAppDefault() *firebase.App {
	// [START initialize_app_default_golang]
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Println("error initializing app: %v\n", err)
	}
	// [END initialize_app_default_golang]

	return app
}

//InitializeServiceAccountID init fcm
func InitializeServiceAccountID() *firebase.App {
	// [START initialize_sdk_with_service_account_id]
	opt := option.WithCredentialsFile("./runex-f341b-firebase-adminsdk-o6pr0-0757998bf7.json")
	conf := &firebase.Config{
		ServiceAccountID: "firebase-adminsdk-o6pr0@runex-f341b.iam.gserviceaccount.com",
		ProjectID:        "runex-f341b",
	}
	app, err := firebase.NewApp(context.Background(), conf, opt)
	if err != nil {
		log.Println("error initializing app: %@ \n", err.Error())
	}
	// [END initialize_sdk_with_service_account_id]
	return app
}

/*func accessServicesSingleApp() (*auth.Client, error) {
	// [START access_services_single_app_golang]
	// Initialize default app
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Access auth service from the default app
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	// [END access_services_single_app_golang]

	return client, err
}

func accessServicesMultipleApp() (*auth.Client, error) {
	// [START access_services_multiple_app_golang]
	// Initialize the default app
	defaultApp, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Initialize another app with a different config
	opt := option.WithCredentialsFile("../runex-f341b-firebase-adminsdk-o6pr0-0757998bf7.json")
	otherApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Access Auth service from default app
	defaultClient, err := defaultApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	// Access auth service from other app
	otherClient, err := otherApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	// [END access_services_multiple_app_golang]
	// Avoid unused
	_ = defaultClient
	return otherClient, nil
}*/
