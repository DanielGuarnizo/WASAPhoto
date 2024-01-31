<script>
export default {
    data: function() {
        return {
            errormsg: null,
            searchUsername: "",
            username: localStorage.getItem("username") || "DefaultUsername",
            token: localStorage.getItem("token")
        }
    },
    methods: {
        async getProfile(u) {
            if (u == "") {
                this.errormsg = "in such a way to search an user the Username cannot be a empty string"
            } else {
                try {
                    //let response = await this.$axios.get(`/users/${this.token}/profile?username=${u}`)
                    localStorage.setItem("searhUsername", u)
                    this.$router.push(`/users/${this.token}/profile?username=${u}`)
                } catch {

                }
            }
        },
    },
    mounted() {
        console.log(' Stream Component is mounted!');
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
                <button class="btn btn-success" type="button" @click="getProfile(searchUsername)">Search </button>
            </div>
        </div>

        <div>
            <button @click="getProfile(this.username)" >
                     <p>{{ this.username }} profile</p>
            </button>
        </div>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
</template>

<style>
</style>