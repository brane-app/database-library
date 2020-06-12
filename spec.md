# Monkebase Spec

This is a spec that should be followed for any database provider that iMonke or `monke-server` interface with

## User things

##### `WriteUser(user User) (err error)`

Take some User `user` and write it to the database. Since `user.ID` is the user's unique key, this should update if that ID already exists, or write a new user if not

##### `ReadSingleUser(ID string) (user User, exists bool, err error)`

Read a single User of ID `ID` and return it, if it exists. `exists` should specify whether or not anything was read.

## Content things

##### `WriteContent(content Content) (err error)`

Write or update a single Content `content` in a database. This should work the same way as WriteUser

##### `ReadSingleContent(ID string) (content Content, exists bool, err error)`

Read a single User of ID `ID` and return it, if it exists. `exists` should specify whether or not anything was read.

##### `ReadManyContent(index, limit int) (content []Content, size int, err error)`

Read up to `limit` Content, starting at `index`. Out of bounds reads will not populate the returned `content`. Size will be the number of items read into the returned array

Contents are ordered by publish timestamp

##### `ReadAuthorContent(ID string, index, limit int) (content []Content, size, err error)`

Read up to `limit` Content, starting at `index` which has been authored by a user of ID `ID`

This should work work in the same way as ReadManyContent, but for a single user

## Auth things

Any passwords stored in auth must be encrypted and salted, preferably using `bcrypt`

Active tokens must expire (usually after 24h)
Active tokens must be stored in memory (not written to a database)

For any user, they must have only one active secret and and token

##### `ReadTokenStat(token string) (owner string, valid bool, err error)`

Read the stat of some token `token`. Returns whether or not this token is `valid` and, if it is, it's owner

If the token is expired or does not exist, it is not valid and will have no owner

##### `CreateToken(ID string) (token string, expires int64, err error)`

Create a token for a user of ID `ID`, invalidating any other token that this user may have active. This token must have an expiry time, returned as a UNIX timestamp `expires`

##### `CreateSecret(ID string) (secret string, err error)`

Create a secret for a user of ID `ID`, invalidating any other secret that this user may have active.

##### `RevokeToken(token string) (err error)`

Revoke a token `token`. If that token is not active, do nothing

##### `RevokeSecret(secret string) (err error)`

Revoke a secret `secret`. If that secret is not active, do nothing

##### `RevokeTokenOf(ID string) (err error)`

Revoke a token belonging to user of ID `ID`. If they have no token, or no such user exists, do nothing

##### `RevokeSecretOf(ID string) (err error)`

Revoke a secret belonging to user of ID `ID`. If they have no secret, or no such user exists, do nothing
