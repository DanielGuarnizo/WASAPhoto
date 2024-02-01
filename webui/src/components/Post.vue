<script>
export default {
    data: function() {
        return { 
            errormsg: null, 
        }
    }, 
    methods:{
        loadimage() {

        },
    },
    computed: { 

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
}
</script>

<template>
    <div class="instagram-post">
        <!-- Post content -->
        <div class="post-content">

        <p>Uploaded: {{ postData.uploaded }}</p>
        <img :src="'data:image/png;base64,' + postData.image" alt="Post Image" class="post-image" />
        
        <p>Number of Comments: {{ postData.numberOfComments }}</p>
        <p>Number of Likes: {{ postData.numberOfLikes }}</p>
        </div>
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