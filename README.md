# Skeleton API
This is a REST-API server template with some predefined stuff. Don't waste time on setting things up.

## Usage
- [Install](https://github.com/tmrts/boilr/wiki/Installation) boilr
- Clone this repository
    ```
    $ cd $GOPATH/src/bitbucket.org/gismart/
    $ git clone git@bitbucket.org:gismart/skeleton-api.git
    ```
- Run template registry initialization
   ```
   $ boilr init
   ```
- Save template to the boilr
    ```
    $ boilr template save $GOPATH/src/bitbucket.org/gismart/skeleton-api skeleton
    ```
- Create new project use this template
    ```
    $ boilr template use skeleton $GOPATH/src/bitbucket.org/gismart
    ```
####Note on DataDog APM usage
Change DataDog {{projectName}} in main.go to your project name. To run DataDog locally start DataDog APM Agent. 
