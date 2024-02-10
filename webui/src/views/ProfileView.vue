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

            } catch {

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
            console.log("inside the get user profile")
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
                if (response.status === 204) {
                    this.banisher = true
                } 
            } catch {
            

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

            } catch {

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
            } catch {

            }
        },
        async followUser() {
            try {
                const response = await this.$axios.post(`/users/${this.token}/following`,{
                    searchUsername: this.searchUsername
                },
                {
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': this.token
                    } 
                });
                this.follower = true 
                console.log(response.status);
                console.log("realoading page")
                window.location.reload();
            } catch (error) {
                // Handle errors
                console.error(error);
            }
        },
        async unfollowUser() {
            try {
                const response = await this.$axios.delete(`/users/${this.token}/following/${this.searchUsername}`,{
                    headers : {
                        'Authorization': `${this.token}`
                    }
                })
                this.follower = false
                console.log(response.status)
                console.log("realoading page")
                window.location.reload();
            } catch {

            }
        },
        async isFollower() {
            try {
                const response = await this.$axios.get(`/users/${this.token}/following`, {
                    headers: {
                        Authorization: this.token
                    } 
                })
                const UserList = response.data;
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

            } catch {

            }
        },

        openFileSelector() {
            // Programmatically trigger the file input
            this.$refs.fileInput.click();
        },
        
        async uploadFile(event) {
            // Handle file selection here
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
                console.log("realoading page")
                window.location.reload();
            } catch (error) {
                // Handle errors
                console.error(error);
            }
        },
    },
    mounted() {
        console.log("user profile Component is mounted")
        this.checkIfBanished()
        this.checkIfBanisher()
        this.isFollower()
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
            <button v-if="!banisher" class="btn btn-success" type="button" @click="banUser()" > Band User </button>
            <button v-else class="btn btn-success" type="button" @click="unBanUser()" > Unband User </button>
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
                    <button @click="getFollowers">Number of Followers {{ this.profile.userFollowers }}</button>
                    <button @click="getFollowing">Number of Followings {{ this.profile.userFollowing }}</button>
                    
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


</style>