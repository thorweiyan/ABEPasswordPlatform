// Vue.use(ElementUI);
// var encrypted = CryptoJS.SHA256('123');
// console.log(encrypted)
// console.log(encrypted.toString())

var app = new Vue({
    el: '#app',
    data: {
        info:{
            header: 'ABE Password Platform',
            footer: '@Copyright by Blockchain Lab @Fudan University'
        },
        user:{
            name:'',
            pass:'',
            passpass:'',
            SFZ:'',
            SJ:'',
            YX:'',
            tip:'',
            policy:{
                tip:'',
                pass:''
            },
            attr:[],
            a:[],
        },
        newAttr:{
            show:false,
            value:'',
        },
    },
    methods:{
        Register(){
            let passHash = CryptoJS.SHA256(this.user.pass).toString()
            let attrs = []
            let tips = []
            attrs.push(`UN:${this.user.name}`)
            attrs.push(`SFZ:${this.user.SFZ}`)
            attrs.push(`SJ:${this.user.SJ}`)
            attrs.push(`YX:${this.user.YX}`)
            // this.user.policy.tip = `${attrs.join(' AND ')}`
            for (let i = 0; i < this.user.attr.length; i++) {
                const e = this.user.attr[i];
                attrs.push(`ZS:${e}`)
                tips.push(`${e}`)
            }
            // this.user.policy.pass = `${attrs.join(' AND ')}`

            // let _Post = {
            //     UserName:`UN:${this.user.name}`,
            //     UserPasswordHash:passHash,
            //     ChangePasswordPolicy:`(${this.user.policy.pass})`,
            //     GetTipPolicy:`(${this.user.policy.tip})`,
            //     GetTipMessage:tips.join(' & '),
            //     UserAttributes:attrs.join(','),
            // }

            let _Post = {
                UserName:`UN:${this.user.name}`,
                UserPasswordHash:passHash,
                ChangePasswordPolicy:`(${this.user.policy.pass})`,
                GetTipPolicy:`(${this.user.policy.tip})`,
                GetTipMessage:tips.join(' & '),
                UserAttributes:attrs.join(','),
            }

            // data post here
            console.log(_Post)
            let _this = this

            axios.post('/signup', _Post)
            .then(function (response) {
                console.log(response);
                _this.$notify.info({
                    title: 'success',
                    message: JSON.stringify(response,null,4),
                    duration: 0
                });
            })
            .catch(function (error) {
                console.log(error);
                _this.$notify.error({
                    title: 'error',
                    message: JSON.stringify(error,null,4),
                    duration: 0
                });
            });
        },
        GenPolicy(){
            let attrs = []
            let tips = []
            attrs.push(`UN:${this.user.name}`)
            attrs.push(`SFZ:${this.user.SFZ}`)
            attrs.push(`SJ:${this.user.SJ}`)
            attrs.push(`YX:${this.user.YX}`)
            this.user.policy.tip = `${attrs.join(' AND ')}`
            for (let i = 0; i < this.user.attr.length; i++) {
                const e = this.user.attr[i];
                attrs.push(`ZS:${e}`)
                tips.push(`${e}`)
            }
            this.user.policy.pass = `${attrs.join(' AND ')}`
        },
        ChangePass(){
            let passHash = CryptoJS.SHA256(this.user.pass).toString()
            let attrs = []
            attrs.push(`UN:${this.user.name}`)
            attrs.push(`SFZ:${this.user.SFZ}`)
            attrs.push(`SJ:${this.user.SJ}`)
            attrs.push(`YX:${this.user.YX}`)
            for (let i = 0; i < this.user.attr.length; i++) {
                const e = this.user.attr[i];
                attrs.push(`ZS:${e}`)
            }

            let _Post = {
                UserName:`UN:${this.user.name}`,
                UserPasswordHash:passHash,
                UserAttributes:attrs.join(','),
            }

            // data post here
            console.log(_Post)
            let _this = this
            axios.post('/changepassword', _Post)
            .then(function (response) {
                console.log(response);
                _this.$notify.info({
                    title: 'success',
                    message: JSON.stringify(response,null,4),
                    duration: 0
                });
            })
            .catch(function (error) {
                console.log(error);
                _this.$notify.error({
                    title: 'error',
                    message: JSON.stringify(error,null,4),
                    duration: 0
                });
            });
        },
        GetTip(){
            let attrs = []
            attrs.push(`UN:${this.user.name}`)
            attrs.push(`SFZ:${this.user.SFZ}`)
            attrs.push(`SJ:${this.user.SJ}`)
            attrs.push(`YX:${this.user.YX}`)
            for (let i = 0; i < this.user.attr.length; i++) {
                const e = this.user.attr[i];
                attrs.push(`ZS:${e}`)
            }

            let _Post = {
                UserName:`UN:${this.user.name}`,
                UserAttributes:attrs.join(','),
            }

            // data post here
            console.log(_Post)
            let _this = this
            axios.post('/gettip', _Post)
            .then(function (response) {
                console.log(response);
                _this.$notify.info({
                    title: 'success',
                    message: JSON.stringify(response,null,4),
                    duration: 0
                });
            })
            .catch(function (error) {
                console.log(error);
                _this.$notify.error({
                    title: 'error',
                    message: JSON.stringify(error,null,4),
                    duration: 0
                });
            });
        },
        Login(){
            let passHash = CryptoJS.SHA256(this.user.pass).toString()

            let _Post = {
                UserName:`UN:${this.user.name}`,
                UserPasswordHash:passHash,
            }

            // data post here
            console.log(_Post)
            let _this = this
            axios.post('/login', _Post)
            .then(function (response) {
                console.log(response);
                _this.$notify.info({
                    title: 'success',
                    message: JSON.stringify(response,null,4),
                    duration: 0
                });
            })
            .catch(function (error) {
                console.log(error);
                _this.$notify.error({
                    title: 'error',
                    message: JSON.stringify(error,null,4),
                    duration: 0
                });
            });
        },

        changeTab(tab){
            console.log(tab.label)
            this.clearContent()
        },
        clearContent(){
            this.user = {
                name:'',
                pass:'',
                passpass:'',
                SFZ:'',
                SJ:'',
                YX:'',
                tip:'',
                policy:{
                    tip:'',
                    pass:''
                },
                attr:[],
                a:[],
            }
        },

        delAttr(attr) {
            this.user.attr.splice(this.user.attr.indexOf(attr), 1);
        },

        showInput() {
            if (this.user.attr.length>=2) {
                this.$notify({
                    title: '警告',
                    message: '自设属性不能超过 2 个',
                    type: 'warning',
                    duration: 0
                });
                return
            }
            this.newAttr.show = true;
            this.$nextTick(_ => {
                this.$refs.save.$refs.input.focus();
            });
        },

        handleInputConfirm() {
            let v = this.newAttr.value;
            if (v&&this.user.attr.indexOf(v)==-1) {
                this.user.attr.push(v);
            }
            this.newAttr.show = false;
            this.newAttr.value = '';
        },

        confirmPass(){
            if (this.user.pass!=''&&
                this.user.passpass!=''&&
                !(this.user.pass===this.user.passpass)) {
                _this.$notify.error({
                    title: '密码不一致',
                    message: '两次输入的密码不同, 请重新输入.'
                });
            }
        },
    },
    computed:{
        // 仅读取
        canRegister() {
            let ok = true
            ok = ok && this.user.name!=''
            ok = ok && this.user.pass!=''
            ok = ok && this.user.passpass!=''
            ok = ok && this.user.SFZ!=''
            ok = ok && this.user.SJ!=''
            ok = ok && this.user.YX!=''
            ok = ok && this.user.policy.pass!=''
            ok = ok && this.user.policy.tip!=''

            ok = ok && (this.user.pass===this.user.passpass)

            return true
            // return ok
        },
        canGetTip() {
            let ok = true
            ok = ok && this.user.name!=''
            ok = ok && this.user.SFZ!=''
            ok = ok && this.user.SJ!=''
            ok = ok && this.user.YX!=''
            return true
            // return ok
        },
        canChangePass() {
            let ok = true
            ok = ok && this.user.name!=''
            ok = ok && this.user.pass!=''
            ok = ok && this.user.SFZ!=''
            ok = ok && this.user.SJ!=''
            ok = ok && this.user.YX!=''
            return true
            // return ok
        },
        canLogin() {
            let ok = true
            ok = ok && this.user.name!=''
            ok = ok && this.user.pass!=''

            return true
            // return ok
        },
    }
  })