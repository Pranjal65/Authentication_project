

package main 

import(
	"encoding/json"
	"fmt"
	// "log"
	"net/http"
	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
)

var(
	redisClient *redis.Client
)

type User struct{
	UserName string `json:username`
	Password string `json:password`
}

func init(){
	redisClient=redis.NewClient(&redis.Options{
		Addr:"localhost:6379",
		Password: "",
		DB:0,
	})
}
func HandleRegistration(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
        var user User
		err:=json.NewDecoder(r.Body).Decode(&user);
		if(err!=nil){
			http.Error(w,"Error during decoding",http.StatusConflict);
			return;
		}
		ifexist,err:=redisClient.Exists(redisClient.Context(),user.UserName).Result()
		if(err!=nil){
            http.Error(w,"Error during checking user",http.StatusConflict);
			return;
		}
		if(ifexist==1){
			http.Error(w,"User laready exists",http.StatusConflict)
			return;
		}
		err=redisClient.Set(redisClient.Context(),user.UserName,user.Password,0).Err()
		if(err!=nil){
			http.Error(w,"Error during registration",http.StatusConflict);
			return;
		}
		//w.WriteHeader(http.StatusCreated)
		// 	fmt.Fprintf(w, "User registered successfully")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w,"User Registered Successfully");
		
		
}
func HandleLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
        var user User
        err:=json.NewDecoder(r.Body).Decode(&user)
		if(err!=nil){
			http.Error(w,"error while decoding",http.StatusConflict)
			return;
		}

		exist,err:=redisClient.Exists(redisClient.Context(),user.UserName).Result()
		if(err!=nil){
            http.Error(w,"error while searching user",http.StatusConflict)
			return
		}
		if(exist!=1){
			http.Error(w,"User does not exist",http.StatusConflict);
			return;
		}

		password,err:=redisClient.Get(redisClient.Context(),user.UserName).Result();
		//fmt.Fprintf(w,password)
		if(err!=nil){
			http.Error(w,"User does not exist",http.StatusUnauthorized);
			return;
		}

		if(password!=user.Password){
                http.Error(w,"Invalid Password",http.StatusUnauthorized);
				return;
		}
		w.WriteHeader(http.StatusOK);
		fmt.Fprintf(w,"Login Successful!");


}

func main(){
	router:=httprouter.New()
	router.POST("/register",HandleRegistration)
	router.POST("/login",HandleLogin)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080",router);

}