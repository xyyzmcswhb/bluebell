package mysql


func init(){
	dbcfg := settings.MySQLConfig{
		Host: "127.0.0.1"
		Port: 3306
		User: "root"
		Password: "123456"
		Dbname: "bluebell"
		Max_connections: 200
		Max_idle_col: 50
	}
	err := Init(&dbcfg)
	if err != nil{
		panic(err)
	}
}

fun TestCreatePost(t *testing.T){
	post := models.Post{
		ID:10,        
   		AuthorID:123,   
		CommunityID:1,
		Title:"test",      
		Content:"just a test"    
	}
	err := CreatePost(&post) 
	if err != nil{
		t.Fatalf("CreatePost insert record into mysql failed, err:%v\n",err)
	}
	t.Logf("CreatePost insert record into mysql success")
}