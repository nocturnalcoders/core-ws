Kebutuhan Entity DB
    - Users
        * Id : int
        * name : varchar
        * occupation : varchar
        * email : varchar
        * password_hash : varchar
        * avatar_file_name : varchar
        * role : varchar
        * token : varchar
        * created_at : datetime
        * updated_at : datetime
    
    - Campaigns
        * Id : int
        * user_id : int
        * name : varchar
        * short_description : varchar
        * description : text
        * goal_amount : int //livecount
        * current_amount : int //currentspacecount
        * perks : text
        * backer_count : text //spacecount
        * slug : varchar
        * created_at : datetime
        * updated_at : datetime

    - Campaign Images
        * id : int
        * campaign_id : int 
        * file_name : varchar
        * is_primary : boolean (tinyint)
        * created_at : datetime
        * updated_at : datetime
        
    - Transactions
        * id : int
        * campaign_id : int
        * user_id : int
        * amount : int
        * status : varchar
        * code : varchar
        * created_at : datetime
        * updated_at : datetime



Step Logic 

input -> Handler Mapping Input ke Struct -> Serivce Mapping ke Struct User (Memanggil Business Process) -> Repository save struct User ke Db (Pemanggilan Database) -> Database
 
Gambaran Proses
User input lalu , handler menangkap data user


INPUT
HANDLER : Mapping input dari user -> input struct
SERVICE : Melakukan mapping dari struct input ke struct user
REPOSITORY
DB

//git tutor
cara memasukan ke git 
1. Git Init terlebih dahulu
2. git add . 
3. git -m "commit nya"
4. git remote add origin https://github.com/nocturnalcoders/BackendEkost.git

//Tentang Middleware
//1. Ambil Nilai Header Autorization 
//Saat user mengirim req ke endpoint yang membutuhkan autorization
//:Bearer tokentokentoken
//2. Dari header authorization, kita ambil nilai token saja
//3. Kita valdiasi token -> pakai service validate token
//4. Menentukan valid / tidak valid
//5. Lalu diambil nilai user_id
//6. ambil user dari db berdasarkan user_id lewat service -> butuh user service
//7. Jika User ada , kita set user context(Sebuah tempat untuk simpan nilai untuk di get / diset di tempat lain)