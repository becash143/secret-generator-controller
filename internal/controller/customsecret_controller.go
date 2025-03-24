package controller

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	appv1 "github.com/becash143/secret-generator-controller/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// CustomSecretReconciler reconciles a CustomSecret object
type CustomSecretReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=app.mydomain.com,resources=customsecrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=app.mydomain.com,resources=customsecrets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=app.mydomain.com,resources=customsecrets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop
func (r *CustomSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the CustomSecret object
	customSecret := &appv1.CustomSecret{}
	if err := r.Get(ctx, req.NamespacedName, customSecret); err != nil {
		log.Error(err, "unable to fetch CustomSecret")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Generate the secret based on the CustomSecret spec
	secretName := fmt.Sprintf("%s-secret", customSecret.Name)
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: customSecret.Namespace,
		},
		Data: make(map[string][]byte),
	}

	// Generate the secret based on SecretType
	switch customSecret.Spec.SecretType {
	case "basic-auth":
		// Generate basic-auth credentials (username: admin, password: random)
		username := customSecret.Spec.Username
		if username == "" {
			username = "admin"
		}

		password := generateRandomString(customSecret.Spec.PasswordLength)

		secret.Data["username"] = []byte(username)
		secret.Data["password"] = []byte(password)

	case "jwt":
		// Generate JWT secret (you can modify this to generate a real JWT token)
		secret.Data["jwt"] = []byte("dummy-jwt-token")

	default:
		log.Error(fmt.Errorf("invalid SecretType"), "Invalid SecretType", "SecretType", customSecret.Spec.SecretType)
		return ctrl.Result{}, fmt.Errorf("invalid SecretType")
	}

	// Create the secret in the cluster if it doesn't exist
	if err := r.Create(ctx, secret); err != nil {
		log.Error(err, "unable to create secret", "secret", secretName)
		return ctrl.Result{}, err
	}

	// Update the CustomSecret status with the secret name and the current time
	customSecret.Status.SecretName = secret.Name
	customSecret.Status.LastUpdated = time.Now().Format(time.RFC3339)

	// Update the CustomSecret object status
	if err := r.Status().Update(ctx, customSecret); err != nil {
		log.Error(err, "unable to update CustomSecret status")
		return ctrl.Result{}, err
	}

	// Schedule a rotation of the secret if a rotation period is specified
	if customSecret.Spec.RotationPeriod != "" {
		rotationPeriod, err := time.ParseDuration(customSecret.Spec.RotationPeriod)
		if err != nil {
			log.Error(err, "unable to parse RotationPeriod")
			return ctrl.Result{}, err
		}

		// Requeue the reconcile loop to trigger a secret rotation after the rotation period
		return ctrl.Result{
			RequeueAfter: rotationPeriod,
		}, nil
	}

	return ctrl.Result{}, nil
}

// generateRandomString generates a random string of a given length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	_, err := rand.Read(result)
	if err != nil {
		panic("failed to generate random string")
	}

	for i := range result {
		result[i] = charset[result[i]%byte(len(charset))]
	}

	return string(result)
}

// SetupWithManager sets up the controller with the Manager.
func (r *CustomSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1.CustomSecret{}).
		Named("customsecret").
		Complete(r)
}
