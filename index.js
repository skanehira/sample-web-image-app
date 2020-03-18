const app = new Vue({
  el: '#app',
  data: {
    header: ["image", "file", "size"],
    files: [],
  },
  methods: {
    getFiles() {
      axios.get("/files").then(response => {
        const files = []
        for (const data of response.data) {
          data.content = `data:image/png;base64,${data.content}`
          files.push(data)
        }

        this.files = files
      })
    },
    upload() {
      const file = document.getElementById("file").files[0]
      const form = new FormData()
      form.append("file", file)

      axios.post("/files", form).then(() => {
        alert("upload successed")
      }).catch(error => {
        console.log(error);
      })
    },
  },
  mounted(){
    this.getFiles()
  }
})


