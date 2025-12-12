package user

// PasswordHasher defines the interface for password hashing operations.
// This abstraction allows the usecase to be independent of the hashing implementation.
type PasswordHasher interface {
	// Hash generates a hash from the given password
	Hash(password string) (string, error)

	// Compare checks if the password matches the hash
	Compare(hashedPassword, password string) error
}
