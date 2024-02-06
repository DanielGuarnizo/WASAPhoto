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
                    let response = await this.$axios.put(`/comments`, {
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
        async uncommentPhoto() {},
    },
    mounted() {
        this.getComments()
    }
}
</script>

<template>
    <div class="input-group mb-3">
        <input type="text" id="comment" v-model="inputComment" class="form-control"
            placeholder="Insert a comment" aria-label="Recipient's username"
            aria-describedby="basic-addon2">
        <div class="input-group-append">
            <button class="btn btn-success" type="button" @click="this.commentPhoto()">Comment</button>
        </div>
    </div>
    <Comment v-for="comment in CommentList" :commentData="comment" />

</template>

<style>

</style>