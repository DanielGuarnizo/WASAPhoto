<script>
import Comment from '@/components/Comment.vue'
export default {
    components: {
        Comment
    },
    data: function() {
        return {
            errormsg: null,
            token: localStorage.getItem("token"),
            usernameLogin: localStorage.getItem("usernameLogin"),
            serachUsername: localStorage.getItem("serachUsername"),
            postid: localStorage.getItem("postid"),
            CommentList: [],
            inputComment: '',
            comment: {
                post_id: '',
                comment_id: '',
                commenter: '',
                user_id: '',
                body: ''
            }
        }
    },
    methods: {
        async getComments() {
            
            try {
                let response = await this.$axios.get(`/users/${this.token}/posts/${this.postid}/comments`, {
                    headers: {
                        Authorization: this.token
                    }
                })
                this.CommentList = response.data
                console.log(response.data)
            } catch (e) {
                console.log(e)
            }
        },
        async commentPhoto(){
            if (this.inputComment){
                try {
                    let response = await this.$axios.post(`/users/${this.token}/posts/${this.postid}/comments`, {
                        post_id: this.postid,
                        comment_id: '',
                        commenter: localStorage.getItem("commenter"),
                        user_id: this.token,
                        body: this.inputComment
                    },{
                        headers: {
                            'Content-Type': 'application/json',
                            'Authorization': `${this.token}`
                        }
                    });
                    console.log(response.data)
                    this.inputComment = '';
                    this.getComments()
                    if (response.status === 204) {
                        return 
                    }
                } catch (e) {
                    console.log(e)
                }
            } else {
                this.errormsg= "before to comment a post insert a text "
            }
        },
        goBack() {
            this.$router.go(-1);
        }
    },
    mounted() {
        this.getComments()
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
                    <button @click="goBack">
                        <!-- <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg> -->
                        go back to the post
                    </button>
                </li>
            </ul>

    
        </div>
    </nav>
    <div class="input-group mb-3">
        <input type="text" id="comment" v-model="inputComment" class="form-control"
            placeholder="Insert a comment" aria-label="Recipient's username"
            aria-describedby="basic-addon2">
        <div class="input-group-append">
            <button class="btn btn-success" type="button" @click="this.commentPhoto()">Comment</button>
        </div>
    </div>
    <Comment v-for="comment in CommentList" :key="comment.comment_id" v-bind="comment" :commentData="comment" />


</template>

<style>

</style>