<script>
export default {
    data: function() {
        return {
            errormsg:false,
            baned:false,
            exits: false,
            baned_by: true,
            searchUsername: localStorage.getItem("searhUsername") || "DefaultUsername",
            username: localStorage.getItem("username"),
            token: localStorage.getItem("token")
        }
    }, 
    methods: {
        async checkBan() {
            // verify if the user with userid can is baned or not from the user with username
            try {
                let response = await this.$axios.get(`/users/${this.token}/bans/${this.searchUsername}`)
                const UserList = response.data;

                // Check is the username is in the list 
                if (UserList === null) {
                    return 
                }
                baned = UserList.includes(username);
                if(baned) {
                    return 
                } else {
                    getUserProfile()
                }
                
            } catch (e) {
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

        },

    },
    mounted() {
        console.log("user profile Component is mounted")
        this.checkBan()
    }
}
</script>

<template>
    <Baned v-if="baned" :msg="this.username"></Baned>
    <div v-else id="personal information">
        <div>
            <p>This is {{ this.username }} profile </p>
        </div>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
</template>

<style>
</style>