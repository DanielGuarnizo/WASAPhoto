<script>
import Post from '@/components/Post.vue';
export default {
    components: {
        Post
    },
    data: function() {
        return {
            errormsg: null,
            searchUsername: "",
            usernameLogin: localStorage.getItem("usernameLogin") || "DefaultUsername",
            token: localStorage.getItem("token"),
            Stream: {
                photos: []
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
        async getMyStream() {
            try {
                let response = await this.$axios.get(`/users/${this.token}/stream`)
                // console.log(response)
                this.Stream = response.data
                // Assuming this.Stream.photos is an array of posts
                this.Stream.photos.sort((a, b) => {
                // Compare the timestamps or dates of the posts in reverse order
                return new Date(b.uploaded) - new Date(a.uploaded);
                });

                // Assign the sorted posts to this.Posts
                this.Posts = this.Stream.photos;
               
                console.log(`this are the post of the Stream: ${this.Posts}`) 
            } catch (e) {
            }
        },
        async getProfile(u) {
            if (u == "") {
                this.errormsg = "in such a way to search an user the Username cannot be a empty string"
            } else {
                try {
                    //let response = await this.$axios.get(`/users/${this.token}/profile?username=${u}`)
                    localStorage.setItem("searchUsername", u)
                    this.$router.push(`/users/${this.token}/profile?username=${u}`)
                } catch {

                }
            }
        },
        
    },
    mounted() {
        console.log(' Stream Component is mounted!');
        this.getMyStream()
    },
    created() {
        
    },
}
</script>

<template>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">

        <div class="input-group mb-3">
            <input 
                type="text" 
                id="username" 
                v-model="searchUsername" 
                class="form-control"
                placeholder="Who you want to search" 
                aria-label="Recipient's username"
                aria-describedby="basic-addon2">
            <div class="input-group-append">
                <button class="btn btn-success" type="button" @click="getProfile(this.searchUsername)">Search </button>
            </div>
        </div>

        <div>
            <button @click="getProfile(this.usernameLogin)" >
                     <p>{{ this.usernameLogin }} profile</p>
            </button>
        </div>
    </div>
    <div class="header-stream">

        <Post v-for="post in this.Posts" :postData="post" />
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
</template>

<style>
.header-stream {
  border: 1px solid #ddd;
  margin: 10px;
  padding: 10px;
  background-color: #fff;
  /* Add other styling for the post container */
}
</style>