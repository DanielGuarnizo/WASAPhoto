<script>
export default {
    data: function() {
        return {
            iscommenter: null,
            usernameLogin: localStorage.getItem("usernameLogin"),
            token: localStorage.getItem("token"),
        }
    },
    props: {
        commentData: {
            type: Object,
            required: true,
            dafault: () => ({
                post_id: '',
                comment_id: '',
                commenter: '',
                user_id: '',
                body: ''
            }),
        }
    },
    methods: {
        async deleteComment() {
            try {
                let response = await this.$axios.delete(`/users/${this.token}/posts/${this.commentData.post_id}/comments/${this.commentData.comment_id}`, {
                    headers: {
                        Authorization: this.token
                    }
                });
                console.log("realoading page")
                window.location.reload();
                // this.$router.push(`/users/${this.token}/profile/${this.commentData.post_id}`)
                if (response.status === 204) {
                    
                }
            } catch (e) {
                console.log(e)
            }
        },
        async isCommenter() {
            this.iscommenter = (this.usernameLogin === this.commentData.commenter)
        }
    },
    mounted() {

        this.isCommenter()
    },
    computed: {

    }

}
</script>

<template>
    <div class="comment-container">
        <div>
            <button v-if="this.iscommenter" class="btn btn-success" type="button" @click="this.deleteComment()">Delete Comment</button>
        </div>
        <div>
            <p class="comment-author">{{ commentData.commenter }}</p>
            <p class="comment-text">{{ commentData.body }}</p>
            
        </div>
    </div>
</template>

<style scoped>
.comment-container {
  border: 1px solid #ddd;
  margin: 10px;
  padding: 10px;
  background-color: #fff;
  /* Additional styling for the comment container */
}

.comment-author {
  margin: 0 0 5px 0;
  font-weight: bold;
  /* Additional styling for the comment author */
}

.comment-text {
  margin: 0;
  /* Additional styling for the comment text */
}
</style>