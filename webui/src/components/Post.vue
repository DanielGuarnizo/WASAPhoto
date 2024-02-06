<script>

export default {
    data: function() {
        return { 
            errormsg: null, 
            owner: localStorage.getItem("owner"),
            searchUsername: localStorage.getItem("searchUsername"),
            usernameLogin: localStorage.getItem("usernameLogin"),
            token: localStorage.getItem("token"), 
            liked: null,
        }
    }, 
    props: {
        postData: {
            type: Object,
            required: true,
            default: () => ({ // Optional: Provide a default value if needed
                user_id: '',
                post_id: '',
                uploaded: '',
                image: '',
                comments: [],
                numberOfComments: 0,
                likes: [],
                numberOfLikes: 0
            }),
        }
    },
    methods:{
        async deletePhoto() {
            console.log("inside the deletePhoto method")
            try {
                let response = await this.$axios.delete(`/users/${this.token}/posts/${this.postData.post_id}`, {
                    headers: {
                        'Authorization': `${this.token}`
                    }
                })
                console.log(response.status)
                console.log("realoading page")
                window.location.reload();
                // this.$router.push(`/users/${this.token}/profile?username=${this.usernameLogin}`)
                // this.$router.push(`/users/${this.token}/stream`)
            } catch {

            }
        },
        async likePhoto() {
            try {
                let response = await this.$axios.put(`/users/${this.token}/posts/${this.postData.post_id}/likes`,{
                    headers : {
                        'Authorization': `${this.token}`
                    }
                })
                console.log(response.data)
                this.liked = true
                console.log("realoading page")
                window.location.reload();
                
            } catch (e) {
                console.error(e);
            }
        },
        async unlikePhoto() {
            try {
                let response = await this.$axios.delete(`/users/${this.token}/posts/${this.postData.post_id}/likes`,{
                    headers : {
                        'Authorization': `${this.token}`
                    }
                })
                if(response.status === 204) {
                    this.liked = false 
                    console.log("realoading page")
                    window.location.reload();
                    return
                }
            } catch (e){
                console.error(e);
            }
        }, 
        async isLiked() {
            try {
                let response = await this.$axios.get(`/users/${this.token}/posts/${this.postData.post_id}/likes`,{
                    headers : {
                        'Authorization': `${this.token}`
                    }
                })
                const UserList = response.data;
                //console.log(`this is the response.data of post ${this.postData.post_id}: ${response.data}, and the response of check the list is : ${UserList.includes(this.usernameLogin)}`)
                if (UserList == null) {
                    this.liked = false
                }
                else if (UserList.includes(this.usernameLogin)) {
                    this.liked = true
                } else {
                    this.liked =  false 
                }
            } catch (e) {
                console.error(e);
            }
        },
        seeComments() {
            localStorage.setItem("postid", this.postData.post_id);
            localStorage.setItem("commenter", this.usernameLogin)
            this.$router.push(`/users/${this.token}/profile/${this.postData.post_id}`)
        }
    },
    mounted() {
        this.isLiked()
    },
    computed: { 
        Owner() {
            return this.searchUsername === this.usernameLogin
        },

    },
}
</script>

<template>
    <div class="instagram-post">
        <div v-if="Owner">
            <button class="btn btn-success" @click="this.deletePhoto()"> 
                Delete Photo
            </button>
        </div>
        <!-- Post content -->
        <div class="post-content">

        <p>Uploaded: {{ postData.uploaded }}</p>
        
        
        <img :src="postData.image" alt="Post Image" class="post-image" />
        
        <p>Number of Comments: {{ postData.numberOfComments }}</p>
        <p>Number of Likes: {{ postData.numberOfLikes }}</p>
        </div>

        <div>
            <button v-if="!this.liked" @click="this.likePhoto()">
                Like
            </button>
            <button v-else @click="this.unlikePhoto()">
                unlike
            </button>
        </div>
        <button class="btn btn-success" @click="this.seeComments()"> Comments</button>
    </div>
   

</template>

<style scoped>
.instagram-post {
  border: 1px solid #ddd;
  margin: 10px;
  padding: 10px;
  background-color: #fff;
  /* Add other styling for the post container */
}


</style>