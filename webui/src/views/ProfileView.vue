<script>
export default {
    data: function() {
        return {
            errormsg:null,
            exits: false,
            baned_by: false,
            searchUsername: localStorage.getItem("searchUsername") || "DefaultUsername",
            usernameLogin: localStorage.getItem("usernameLogin"),
            token: localStorage.getItem("token")
        }
    }, 
    methods: {
        async checkBan() {
            // verify if the user with userid can is baned or not from the user with username
            console.log("before the try catch of checkBan")
            try {
                let response = await this.$axios.get(`/users/${this.token}/bans/${this.searchUsername}`)
                const UserList = response.data;
                console.log("insede the try of checkBan")
                console.log(response.data)

                // Check is the username is in the list 
                if (UserList === null) {
                    this.baned_by = false 
                    return 
                }
                
                this.baned_by = UserList.includes(this.usernameLogin);
                console.log("result of the baned operation", this.baned_by)
                if(this.baned_by) {
                    return 
                } else {
                    getUserProfile()
                }
                
            } catch (e) {
                console.log("inside the catch of checkBan")
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
            try {
                let reponses = await this.$axios.get(`/users/${this.token}/profile?username=${this.searchUsername}`)

            } catch {

            }
        },

    },
    mounted() {
        console.log("user profile Component is mounted")
        this.checkBan()
    }
}
</script>

<template>
    <Baned v-if="baned_by" :msg="this.searchUsername"></Baned>
    <div v-else id="personal information">
        <div>
            <p>This is {{ this.searchUsername }} profile </p>
        </div>
    </div>
    <ErrorMsg v-if="this.errormsg" :msg="this.errormsg"></ErrorMsg>
</template>

<style>
</style>