
Собрать image в Docker:
                docker build -t go-webapp .

Варианты запуска 
                docker run --publish 4000:4000  --name test --rm  go-webapp
                docker run --publish 4000:4000  --name test --rm -e ENVDBTYPE="mysql"  go-webapp
                docker run --publish 4000:4000  --name test --rm -e ENVDBTYPE="sqlite3"  go-webapp
                docker run --publish 4000:4000  --name test --rm -e ENVDBTYPE="postgres"  go-webapp
                docker run --publish 4000:4000  --name test --rm -e ENVDBTYPE="sqlite3" -e  ADR_GRPC="172.17.0.2:9090"  go-webapp

Envirements
    ENVDBTYPE:
                mysql
                sqlite
                postgres

Команда для добавления b просмотра инцидента под пользователем joe
                curl--request POST --url http://localhost:4000/addIncident --header "Authorization: Basic am9lOjEyMzQ="
                curl --request GET --url http://localhost:4000/showIncident?id=1 --header "Authorization: Basic am9lOjEyMzQ="
                
                

