<script>
export default {
	data: function() {
		return {
			errormsg: null,
			username: "",
			user: {
				user_id: "",
				username: this.username
			},
		}
	},
	methods: {
		async doLogin() {
			if (this.username == "") {
				this.errormsg = "Username cannot be a empty string"
			} else {
				// we make use of the asyn-await here 
				try {
					let response = await this.$axios.post("/session", { username: this.username})
					this.user.user_id = response.data.user_id
					localStorage.setItem("token", this.user.user_id);
                    localStorage.setItem("username", this.username);
										
					this.$router.push(`/users/${this.user.user_id}/stream`)
				} catch (e){
					if (e.response && e.response.status === 400) {
                        this.errormsg = "Form error, may you write in a incorrect way the username . If you think that this is an error, write an e-mail to us.";
                        this.detailedmsg = null;
                    } else if (e.response && e.response.status === 500) {
                        this.errormsg = "Internal server error occurred. Please try again later.";
                        this.detailedmsg = e.toString();
                    } else {
                        this.errormsg = e.toString();
                        this.detailedmsg = null;
                    }
				}
			}
		},
	},
	mounted() {
		console.log('Login Component has been mounted');
		//this.doLogin()
	}
}
</script>

<template>
    <div
        class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
        <h1 class="h2">Welcome to WASAPhoto</h1>
    </div>
    <div class="input-group mb-3">
        <input type="text" id="username" v-model="username" class="form-control"
            placeholder="Insert a username to log in WASAPhoto." aria-label="Recipient's username"
            aria-describedby="basic-addon2">
        <div class="input-group-append">
            <button class="btn btn-success" type="button" @click="doLogin">Login</button>
        </div>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
</template>

<style>
</style>