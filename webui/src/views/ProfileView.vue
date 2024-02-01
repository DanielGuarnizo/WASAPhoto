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
        }
    }, 
    methods: {
        async refresh() {
            
        },
        async checkIfBanisher() {
            try {
                let response = await this.$axios.get(`/users/${this.token}/bans/${this.usernameLogin}`, {
                    headers: {
                        Authorization: this.token
                    } 
                })
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
                let response = await this.$axios.get(`/users/${this.token}/bans/${this.searchUsername}`, {
                    headers: {
                        Authorization: this.token
                    } 
                })
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
                //console.log("get user profile try")
                let response = await this.$axios.get(`/users/${this.token}/profile?username=${this.searchUsername}`, {
                    headers: {
                        Authorization: this.token
                    }
                })
                this.profile = response.data
                this.Posts = this.profile.photos                
            } catch (e){
                console.log("get user profile catch ")
                if (e.response && e.response.status === 400) {
                    this.errormsg = "Bad Request. invalid user data";
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
        async banUser() {
            try {
                let response = await this.$axios.put(`/users/${this.token}/bans/${this.searchUsername}`, {
                    headers: {
                        Authorization: this.token
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
                    const response = await this.$axios.put(
                        `/users/${this.token}`,
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
                let response = await this.$axios.post()
            } catch {

            }
        },
        async unfollowUser() {},
    },
    mounted() {
        console.log("user profile Component is mounted")
        this.checkIfBanished()
        this.checkIfBanisher()

    }, 
    computed: {
        Owner() {
            return this.searchUsername === this.usernameLogin
        }
    }
}
</script>

<template>
    <div v-if="Owner" class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 ">
        <button class="btn btn-success" type="button" @click="uploadPhoto()"> Upload Photo </button>
        
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
                <button v-if="!Owner" class="btn btn-success" type="button" @click="followUser()" > Follow </button>
                <button v-else class="btn btn-success" type="button" @click="unfollowUser()" > UnFollow </button>
            </div>
            <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 ">
                <p>Number of Posts</p>
                <p>Number of Followers</p>
                <p>Number of Followings</p>
                
            </div>
            <Post v-for="post in Posts" :postData="post" />
        </div>
    </div>
    <ErrorMsg v-if="this.errormsg" :msg="this.errormsg"></ErrorMsg>
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