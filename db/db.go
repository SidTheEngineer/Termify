package db

import (
	"log"
	"strconv"
	"time"

	"github.com/SidTheEngineer/Termify/auth"
	"github.com/SidTheEngineer/Termify/playbackUI"
	"github.com/SidTheEngineer/Termify/util"

	"github.com/boltdb/bolt"
)

const (
	dbName              = "SpotfiyAuth.db"
	accessTokenText     = "accessToken"
	tokenTypeText       = "tokenType"
	tokenScopeText      = "tokenScope"
	refreshTokenText    = "refreshToken"
	tokenExpiresInText  = "tokenExpiresIn"
	timeTokenCachedText = "timeTokenCached"
)

// DB is our bolt database that will hold token information.
var DB *bolt.DB

// Start opens and starts a new BoltDB to be connected to.
func Start() {
	newDb, err := bolt.Open(dbName, 0600, &bolt.Options{Timeout: 5 * time.Second})

	if err != nil {
		log.Fatal(err)
	}

	DB = newDb
	DB.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("auth"))
		return nil
	})
}

// Close closes the db
func Close() {
	DB.Close()
}

// CacheAccessToken caches Spotify access token information in the DB.
// TODO: Move this over from the auth package.
func CacheAccessToken(tx *bolt.Tx, token auth.AccessToken) {
	authBucket := tx.Bucket([]byte("auth"))
	authBucket.Put([]byte(accessTokenText), []byte(token.Token))
	authBucket.Put([]byte(tokenTypeText), []byte(token.Type))
	authBucket.Put([]byte(tokenScopeText), []byte(token.Scope))
	authBucket.Put([]byte(refreshTokenText), []byte(token.RefreshToken))
	authBucket.Put([]byte(tokenExpiresInText), []byte(strconv.Itoa(token.ExpiresIn)))
	authBucket.Put([]byte(timeTokenCachedText), []byte(strconv.FormatInt(int64(time.Now().Unix()), 10)))
}

// IsLoggedIn checks if a user has cached access token information that can be used
// to fetch data from Spotify.
func IsLoggedIn(uiConfig *playbackUI.Config, authConfig *auth.Config) bool {
	loggedIn := true
	DB.Batch(func(tx *bolt.Tx) error {
		authBucket := tx.Bucket([]byte("auth"))
		accessToken := authBucket.Get([]byte(accessTokenText))
		tokenType := authBucket.Get([]byte(tokenTypeText))
		tokenScope := authBucket.Get([]byte(tokenScopeText))
		refreshToken := authBucket.Get([]byte(refreshTokenText))
		expiresIn := authBucket.Get([]byte(tokenExpiresInText))
		timeTokenCached := authBucket.Get([]byte(timeTokenCachedText))

		allFieldsCached := !util.IsNil(
			accessToken,
			tokenType,
			refreshToken,
			expiresIn,
			timeTokenCached,
		) && !util.IsEmpty(
			accessToken,
			tokenType,
			refreshToken,
			expiresIn,
			timeTokenCached,
		)

		if allFieldsCached {
			if auth.TokenIsExpired(string(timeTokenCached), string(expiresIn)) {
				token := auth.FetchSpotifyTokenByRefresh(string(refreshToken))
				CacheAccessToken(tx, token)
				uiConfig.SetAccessToken(token)
				authConfig.SetAccessToken(token)
			} else {
				expireTime, _ := strconv.Atoi(string(expiresIn))

				uiConfig.SetAccessToken(auth.AccessToken{
					Token:        string(accessToken),
					Type:         string(tokenType),
					Scope:        string(tokenScope),
					RefreshToken: string(refreshToken),
					ExpiresIn:    expireTime,
				})
				authConfig.SetAccessToken(auth.AccessToken{
					Token:        string(accessToken),
					Type:         string(tokenType),
					Scope:        string(tokenScope),
					RefreshToken: string(refreshToken),
					ExpiresIn:    expireTime,
				})
			}
			// TODO: Return error that might be generated
			return nil
		}
		loggedIn = false
		return nil
	})
	return loggedIn
}
