<!-- TABLE OF CONTENTS -->

  <summary>Tabla de contenidos</summary>
  <ol>
    <li>
      <a href="#about-the-project">Acerca de este proyecto</a>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisitos">Prerequisitos</a></li>
         <li><a href="#investigando-jira">Investigando Jira</a></li>
        <li><a href="#instalacion">Instacion</a></li>
      </ul>
    </li>
    <li><a href="#todo">TO-DO</a></li>
    <li><a href="#errores-comunes">Errores</a></li>
    <li><a href="#usage">referencias</a></li>
  </ol>


<!-- ABOUT THE PROJECT -->
## About The Project

Configurar y levantar en modo cluster Jira server utilizando docker, creación de ticket por medio de api e identificacion de procesos que ocurren detras de escena.

Espero pueda servir para jugar en el local de cada uno. 

<!-- GETTING STARTED -->
## Getting Started

<!-- prerequisites -->
## Prerequisitos

* docker
* docker pull atlassian/jira-software:9.15.2-ubuntu-jdk17
* docker pull postgres:14.13-bullseye
* Go lang (opcional)

<!-- investigando-jira -->
### Investigando Jira

Para comenzar la idea es investigar un poco Jira e identificar cuales son sus componentes basicos para funcionar

```bash
docker run --name "new-dockerized-jira" -p 8081:8080 atlassian/jira-software:9.15.2-ubuntu-jdk17
```

La instalacion de jira dento del contenedor esta en: 

```bash
/var/atlassian/application-data/jira
```

Los folder y archivos son los siguientes:

```bash
drwxr-x--- 2 jira jira 4096 Oct  8 15:24 analytics-logs
drwxr-x--- 4 jira jira 4096 Oct  8 14:59 caches
drwxr-x--- 3 jira jira 4096 Oct  8 15:25 customisations
drwxr-x--- 5 jira jira 4096 Oct  8 15:24 data
drwxr-x--- 2 jira jira 4096 Oct  8 15:23 database
-rw-r--r-- 1 jira jira  897 Oct  8 15:23 dbconfig.xml
-rw-r--r-- 1 jira jira    1 Oct  8 14:59 docker-app.pid
drwxr-x--- 4 jira jira 4096 Oct  8 15:23 export
drwxr-x--- 4 jira jira 4096 Oct  8 15:23 import
drwxr-x--- 2 jira jira 4096 Oct  8 15:24 keys
drwxr-x--- 4 jira jira 4096 Oct  8 15:26 log
drwxr-x--- 2 jira jira 4096 Oct  8 15:25 logos
drwxr-x--- 2 jira jira 4096 Oct  8 15:23 monitor
drwxr-x--- 6 jira jira 4096 Oct  8 15:23 plugins
drwxr-x--- 3 jira jira 4096 Oct  8 15:28 tmp
```


<!-- instalacion -->
## Instalación


La primera parte es configurar la BD y el nodo 1 de Jira(stannd alone), solo se debe levantar el nodo 1 para poder hacer las configuraciones basicas.

Dentro del yaml de ejemplo se encuentra comentados los nodos 2 y 3 dejar asi hasta las siguiente es estapas.

Se configuraron 6 volumes uno por cada nodo de Jira, uno para la base de datos un volumen compartido entre los nodos

```yaml
volumes:
  postgres-data:
  jira-home-node-01:
  jira-home-node-02:
  jira-home-node-03:
  jira-home-node-04:
  jira-shared-home:
```

Para que los nodos se puedan ver entre ellos es importante habilitar una red compartida

```yanl
networks:
  jira-network:
    driver: bridge
```

Recordar solo levantar el nodo 1 de Jira
```bash
#levantar los contenedores postgresql y nodo 1 de jira
docker compose up 
```

Una vez la base de datos y el nodo 1 de jira esten running, configurar  jira por medio de la url http://localhost:8081/


### Configurar Nodo 2 y otros

Una vez finalizado la instalacion por la web, ingresaremos al contenedor y copiaremos los siguientes archivos al volumen jira-shared-home

```sh
- `data`
- `plugins`
- `logos`
- `import`
- `export`
- `caches`
- `keys`
```

Para poder traspasar las configuraciones entre nodos voy a usar el mismo volume y folder dentro del volume jira-shared-home, para esto creare una carpeta llamada commons  donde estaran todos los archivos que se requiern copiar entre nodos como el archivo dbconfig.xml que debe ser el mismo en todos los nodos

```bash
cp -R data/ plugins/ logos/ import/ export/ caches/ keys/ ../../shared-application-data/jira/.
mkdir /var/atlassian/shared-application-data/jira/common
#copiar la configuracion de base de datos para distribuir en otros nodos
cp dbconfig.xml /var/atlassian/shared-application-data/jira/common
#
touch /var/atlassian/shared-application-data/jira/common/cluster.properties-node-01
touch /var/atlassian/shared-application-data/jira/common/cluster.properties-node-02
touch /var/atlassian/shared-application-data/jira/common/cluster.properties-node-03
touch /var/atlassian/shared-application-data/jira/common/cluster.properties-node-04

```

Agregar lo siguiente a los archivos cluster.properties-node-01 al nodo N

```bash
# This unique ID must match the username and the BalancerMember entry in the Apache config
jira.node.id = jira-pg-node-01 # remplazar por cada nodo
# The location of the shared home directory for all Jira nodes
jira.shared.home = /var/atlassian/shared-application-data/jira
```

Una vez esten todas las configuraciones terminadas levantar el resto de nodos




<!-- todo -->
## TODO

- [ ] add nginx
- [ ] Add back to top links
- [ ] refactor main.go
- [ ] crear funcion que permita crear strings randmos



<!-- errores-comunes -->
## Errores Comunes
### ERROR:  duplicate key value violates unique constraint "cluster_lock_name_idx"

```
2024-10-09 15:56:08 jira-pg-node-01    | 2024-10-09 18:56:08.828 UTC [349] ERROR:  duplicate key value violates unique constraint "cluster_lock_name_idx"
2024-10-09 15:56:08 jira-pg-node-01    | 2024-10-09 18:56:08.828 UTC [349] DETAIL:  Key (lock_name)=(com.atlassian.greenhopper.manager.lexorank.LexoRankCleanupScheduler) already exists.
2024-10-09 15:56:08 jira-pg-node-01    | 2024-10-09 18:56:08.828 UTC [349] STATEMENT:  insert into public.clusterlockstatus (lock_name, update_time, id)
```

 Es normal ver este tipo de errores, esto ocurre cuando los nodos estan  en proceso de inicio
 
### java.rmi.UnknownHostException: Unknown host

ving from queue. Error: java.rmi.UnknownHostException: Unknown host: 69cf866119b7; nested exception is:
    	java.net.UnknownHostException: 69cf866119b7


Uno o varios nodos no estan funcionando correctamente

select * from clusternode;
select * from clusternodeheartbeat;

### Referencias

https://confluence.atlassian.com/adminjiraserver/set-up-a-jira-data-center-cluster-993929600.html

https://developer.atlassian.com/server/jira/platform/configuring-a-jira-cluster/

https://www.docker.com/blog/how-to-use-the-postgres-docker-official-image/

#### Posgresql Docker
https://medium.com/@jewelski/quickly-set-up-a-local-postgres-database-using-docker-5098052a4726


####  JIRA API DEVELOPER

https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-issues/#api-group-issues



### Borrado de ambiente


docker rm jira-node-01 jira-node-02 jira-node-04 jira-pg-node-01  container-pgadmin
docker volume rm deploy_jira-home-node-01 deploy_jira-home-node-02  deploy_jira-home-node-03 deploy_jira-home-node-04 deploy_jira-shared-home deploy_postgres-data
docker network rm deploy_jira-network