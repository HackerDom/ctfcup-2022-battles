package main

import "github.com/golang/glog"

func AddUser(user string, token string) error {
	_, err := db.Exec(`INSERT INTO secrets.secrets.users (name, token) VALUES ($1, $2)`, user, token)
	if err != nil {
		glog.Fatal(err.Error())
		return err
	}
	return nil
}

func AddSecret(owner string, text string) error {
	_, err := db.Exec(`INSERT INTO secrets.secrets.secrets (owner, secret) VALUES ($1, $2)`, owner, text)
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	return nil
}

func GetSecrets(owner string) []string {
	rows, err := db.Queryx(`SELECT secret FROM secrets.secrets.secrets where owner = $1`, owner)
	if err != nil {
		glog.Error(err.Error())
		return []string{}
	}

	var result []string
	for rows.Next() {
		var secret string
		err = rows.Scan(&secret)
		if err != nil {
			glog.Error(err.Error())
			return []string{}
		}

		result = append(result, secret)
		glog.Infof("%+v", secret)
	}

	return result
}

func GetUsers() []User {
	rows, err := db.Queryx(`SELECT  * FROM secrets.secrets.users`)
	if err != nil {
		glog.Error(err.Error())
		return []User{}
	}

	var result []User
	for rows.Next() {
		var user User
		err = rows.StructScan(&user)
		if err != nil {
			glog.Error(err.Error())
			return []User{}
		}

		result = append(result, user)
	}

	return result
}
