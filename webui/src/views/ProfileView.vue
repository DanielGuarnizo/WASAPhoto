<script>
import Post from '@/components/Post.vue';

export default {
    components: {
        Post
    },
    data: function() {
        return {
            errormsg:null,
            exits: false,
            banished: null,
            banisher: null,
            follower: null,
            selectedFile: null,
            base64Image: null,
            formattedDate: '',
            followings: [],
            followers: [],
            searchUsername: localStorage.getItem("searchUsername") || "DefaultUsername",
            usernameLogin: localStorage.getItem("usernameLogin"),
            token: localStorage.getItem("token"),
            newUsername: null,
            profile: {
                user: {
                user_id: '',
                username: '',
                },
                photos: [],
                numberOfPosts: 0,
                userFollowers: 0,
                userFollowing: 0,
            },
            Posts: [],
            Post: {
                user_id: '',
                post_id: '',
                uploaded: '',
                image: '',
                comments: [],
                numberOfComments: 0,
                likes: [],
                numberOfLikes: 0
            }
        }
    }, 
    methods: {
        async checkIfBanisher() {
            try {
                let response = await this.$axios.get(`/users/${this.token}/bans?searchUsername=${this.usernameLogin}`,{
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `${this.token}`
                    }
                });
                const UserList = response.data;
                // Check is the username is in the list if empty then the serachUsername has not baned persons
                if (UserList === null) {
                    this.banisher = false 
                    return 
                }
                this.banisher = UserList.includes(this.searchUsername);
                
                if(this.banisher) {
                    return 
                } else {
                    return
                }

            } catch (e) {
                console.error(e);
            }   
        },
        async checkIfBanished() {
            // verify if the user with userid can is baned or not from the user with username
            console.log("before the try catch of checkIfBanished")
            try {
                let response = await this.$axios.get(`/users/${this.token}/bans?searchUsername=${this.searchUsername}`,
                {
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `${this.token}`
                    }
                });
                const UserList = response.data;
                console.log("before the if inside the try in checkifbanished")
                // Check is the username is in the list if empty then the serachUsername has not baned persons
                if (UserList === null) {
                    this.banished = false 
                    this.getUserProfile();
                    return 
                }
                
                this.banished = UserList.includes(this.usernameLogin);
                console.log("result of the baned operation", this.banished)
                if(this.banished) {
                    return 
                } else {
                    this.getUserProfile();
                }
                
            } catch (e) {
                console.log("inside the catch of checkIfBanished")
                if (e.response && e.response.status === 400) {
                    this.errormsg = "Bad reques. Invalid user data";
                    this.detailedmsg = null;
                } else if (e.response && e.response.status === 500) {
                    this.errormsg = "Internal server error occurred. Please try again later.";
                    this.detailedmsg = e.toString();
                } else {
                    this.errormsg = e.toString();
                    this.detailedmsg = null;
                }
            }
        },
        async getUserProfile() {
            console.log("||| inside the get user profile")
            console.log(this.usernameLogin)
            console.log(this.searchUsername)
            try {
                console.log("get user profile try")
                let response = await this.$axios.get(`/users/${this.token}/profile?username=${this.searchUsername}`, {
                    headers: {
                        Authorization: this.token
                    }
                })
                this.profile = response.data
                this.Posts = this.profile.photos  
                console.log(`this are the post of the profile: ${this.Posts}`) 
            } catch (e) {
                console.log("get user profile catch ")
                if (e.response && e.response.status === 400) {
                    this.errormsg = "Bad Request. invalid user data";
                } else if (e.response && e.response.status === 500) {
                    this.errormsg = "the user you are looking for is not available ";
                } else {
                    console.log("enter in the else of the get user profile")
                    this.errormsg = e.toString();
                }
            }
        },
        async banUser() {
            try {
                console.log(this.searchUsername)
                let response = await this.$axios.post(`/users/${this.token}/bans`,
                {
                    banished: this.searchUsername
                },
                {
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `${this.token}`
                    } 
                })
                console.log(response.data)
                if (response.status === 201) {
                    this.banisher = true
                } 
            } catch (e) {
                console.error(e);
            }
        },
        async unBanUser() {
            try {
                let response = await this.$axios.delete(`/users/${this.token}/bans/${this.searchUsername}`, {
                    headers: {
                        Authorization: this.token
                    } 
                })
                if (response.status === 204) {
                    this.banisher = false
                } 

            } catch (e) {
                console.error(e);
            }
        }, 
        async setMyUserName () {
            try {
                if (this.newUsername) {
                    const response = await this.$axios.put(`/users/${this.token}`,
                        {
                        newUsername: this.newUsername
                        },
                        {
                        headers: {
                            'Content-Type': 'application/json',
                            'Authorization': this.token
                        }
                        }
                    );
                    if (response.status === 200) {
                        this.usernameLogin = this.newUsername
                        this.searchUsername = this.newUsername
                        this.newUsername = null
                    }
                }
            } catch (e) {
                console.error(e);
            }
        },
        async followUser() {
            try {
                const response = await this.$axios.post(`/users/${this.token}/followings`,{
                    searchUsername: this.searchUsername
                },
                {
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': this.token
                    } 
                });
                this.follower = true 
                console.log(response.data);
                console.log(response.status);
                // console.log("realoading page")
                // window.location.reload();
            } catch (error) {
                // Handle errors
                console.error(error);
            }
        },
        async unfollowUser() {
            try {
                const response = await this.$axios.delete(`/users/${this.token}/followings/${this.searchUsername}`,{
                    headers : {
                        'Authorization': `${this.token}`
                    }
                })
                this.follower = false
                console.log(response.status)
                // console.log("realoading page")
                // window.location.reload();
            } catch (e) {
                console.error(e);
            }
        },
        async isFollower() {
            try {
                const response = await this.$axios.get(`/users/${this.token}/followings`, {
                    headers: {
                        Authorization: this.token
                    } 
                })
                const UserList = response.data;
                this.followings = response.data;
                // Check is the username is in the list if empty then the serachUsername has not baned persons
                if (UserList === null) {
                    this.follower = false 
                    return 
                }
                this.follower = UserList.includes(this.searchUsername);
                
                if(this.banisher) {
                    return 
                } else {
                    return
                }

            } catch (e) {
                console.error(e);
            }
        },
        async getFollowers () {
            try {
                const response = await this.$axios.get(`/users/${this.token}/followers`, {
                    headers: {
                        Authorization: this.token
                    } 
                })
                this.followers = response.data
            } catch (e) {
                console.error(e);
            }
        },

        openFileSelector() {
            // Programmatically trigger the file input
            this.$refs.fileInput.click();
        },
        
        async uploadFile(event) {
            // Handle file selection here
            console.log(event.target.files[0])
            if (event.target.files[0]) {
                await new Promise((resolve) => {
                const reader = new FileReader();
                reader.onload = (e) => {
                    this.base64Image = e.target.result;
                    resolve(); // Resolve the promise when the image is loaded
                };
                reader.readAsDataURL(event.target.files[0]);
                });
            } else {
                this.errormsg = "Select an image to upload to your profile";
            }

            // Now that the image is loaded, call the uploadPhoto function
            this.uploadPhoto();
        },
        async uploadPhoto() {
            
            if (!this.base64Image) {
                this.errormsg = "Select an image to upload to your profile";
                return;
            }

            try {
                let response = await this.$axios.post(`/users/${this.token}/posts`, 
                {
                    user_id: this.token,
                    post_id: '',
                    uploaded: this.formattedDate,
                    image: this.base64Image,
                    comments: [],
                    numberOfComments: 0,
                    likes: [],
                    numberOfLikes: 0
                }, 
                {
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `${this.token}`
                    }
                });
                this.errormsg = null;
                this.$refs.fileInput.value = null;
                this.base64Image = null;
                
                
                this.getUserProfile();
                console.log(response.data)
                console.log(response.status)
                console.log(response)
                // console.log("realoading page")
                // window.location.reload();
            } catch (error) {
                // Handle errors
                console.error(error);
            }
        },
        openPopup(stringList) {
            // Get the pop-up container and content elements
            var popupContainer = document.getElementById("popupContainer");
            var popupContent = document.getElementById("popupContent");

            // Simulated list of strings (replace this with your actual data)
            // var stringList = this.followers

            // Build the HTML content for the list of strings
            var listHTML = "<ul>";
            stringList.forEach(function (item) {
                listHTML += "<li>" + item + "</li>";
            });
            listHTML += "</ul>";

            // Set the content of the pop-up
            popupContent.innerHTML = listHTML;

            // Display the pop-up
            popupContainer.style.display = "block";
        },

            // You can also add a function to close the pop-up if needed
        closePopup() {
            var popupContainer = document.getElementById("popupContainer");
            popupContainer.style.display = "none";
        },
        
    },
    mounted() {
        console.log("user profile Component is mounted")
        this.checkIfBanished()
        this.checkIfBanisher()
        this.isFollower()
        this.getFollowers()
        const currentDate = new Date();
        this.formattedDate = currentDate.toISOString().slice(0, 19).replace('T', ' ');
    }, 
    computed: {
        Owner() {
            return this.searchUsername === this.usernameLogin
        },

    }
}
</script>

<template>
    <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
        <div class="position-sticky pt-3 sidebar-sticky">
            <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
                <span>General</span>
            </h6>
            <ul class="nav flex-column">
                <li class="nav-item">
                    <RouterLink to='/users/:userid/stream' class="nav-link">
                        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
                        Home
                    </RouterLink>
                </li>
            </ul>

    
        </div>
    </nav>
    <ErrorMsg v-if="this.errormsg" :msg="this.errormsg"></ErrorMsg>
    <div v-else >
        <div v-if="Owner" class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 ">

            <div>
                <button class="btn btn-success" @click="openFileSelector">Upload Photo</button>
                <input ref="fileInput" type="file" style="display: none;" @change="uploadFile">
            </div>
            
            <input type="text" id="newUsername" v-model="this.newUsername" class="form-control"
            placeholder="Insert the new username" aria-label="Recipient's username"
            aria-describedby="basic-addon2">
    
            <div class="input-group-append">
                <button class="btn btn-success" type="button" @click="setMyUserName()"> Set My User Name </button>
            </div>
    
        </div>
        <div v-else>
            <button v-if="!banisher" class="btn btn-success" type="button" @click="banUser()" > Ban User </button>
            <button v-else class="btn btn-success" type="button" @click="unBanUser()" > Unban User </button>
        </div>
    
    
        <div class="header-profile">
    
            <Baned v-if="banished" :msg="this.searchUsername"></Baned>
            <p v-else-if="banisher"> You band this user then you cannot see his profile</p>
    
            <div v-else id="personal information">
                <div>
    
                    <h2 class="text-center border-bottom">{{ this.searchUsername }} Profile</h2>
                    <div v-if="!Owner">
    
                        <button v-if="!follower" class="btn btn-success" type="button" @click="followUser()" > Follow </button>
                        <button v-else class="btn btn-success" type="button" @click="unfollowUser()" > UnFollow </button>
                    </div>
                </div>
                <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 ">
                    <p>Number of Posts {{ this.profile.numberOfPosts }}</p>
                    <button @click="openPopup(this.followers)">Number of Followers {{ this.profile.userFollowers }}</button>
                    <!-- Pop-up Container -->
                    <div id="popupContainer" class="popup-container" @click="closePopup">
                        <div id="popupContent" class="popup-content">
                            <!-- Content of the pop-up (list of strings) will be inserted here -->
                        </div>
                    
                    </div>
                    <button @click="openPopup(this.followings)">Number of Followings {{ this.profile.userFollowing }}</button>
                    <div id="popupContainer" class="popup-container" @click="closePopup">
                        <div id="popupContent" class="popup-content">
                            <!-- Content of the pop-up (list of strings) will be inserted here -->
                        </div>
                    
                    </div>
                </div>
                <Post v-for="post in Posts" :key="post.post_id" v-bind="post" :postData="post"/>
            </div>
        </div>
    </div>
    
</template>

<style>
.header-profile {
  border: 1px solid #ddd;
  margin: 10px;
  padding: 10px;
  background-color: #fff;
  /* Add other styling for the post container */
}

.popup-container {
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
}

.popup-content {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.3);
}


</style>