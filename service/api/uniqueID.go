package api 

// Function to generate a unique ID
func generateUniqueID() string {
	rand.Seed(time.Now().UnixNano())

	const charset = "0123456789abcdefABCDEF"
	const length = 32

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
