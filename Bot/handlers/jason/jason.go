package jason

import (
	"strings"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/config"
	dbhelper "github.com/ZestHusky/femboy-control/Bot/dbhelpers"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	"github.com/ZestHusky/femboy-control/Bot/logging"
	"github.com/bwmarrin/discordgo"
)

var jasonVideos []string = []string{
	"https://cdn.discordapp.com/attachments/1112127464584511538/1261031430780227634/Y2meta.app_-_Spy_jason_statham_scene_1_online-video-cutter.com_2_cropped.mov?ex=669374c9&is=66922349&hm=5fca00e577e1b96689e19d30fc6bdc8838f157a3a4e1ae7c8a218f620cfd283f&",
	"https://cdn.discordapp.com/attachments/1112127464584511538/1261031430369054730/Y2meta.app_-_Spy_jason_statham_scene_1_online-video-cutter.com.mp4?ex=669374c9&is=66922349&hm=19fdfb1a5f36e7c8ea454aabb3f34cf42358c89044ab6bfea1cd1332030e722c&",
	"https://cdn.discordapp.com/attachments/1112127464584511538/1261031429786042369/The_Most_Annoying_Sound_In_The_World_online-video-cutter.com.mov?ex=669374c9&is=66922349&hm=6e62d521d7ebfc934ecd50ec3d8a41cc52fc86403cefed5d226f49979e466e32&",
	"https://cdn.discordapp.com/attachments/1112127464584511538/1261031429308026900/stethem_daet_pizdi.mov?ex=669374c9&is=66922349&hm=fbaef67b0c48aedc124294aaa86217cbf6d6e7c95438cdb93a843e1509b32632&",
	"https://cdn.discordapp.com/attachments/1112127464584511538/1261031421028601937/DWAR_-_Hor_Hor_with_J.Statham.mov?ex=669374c7&is=66922347&hm=560edb8dde166b37e3b7f11f0f23334d69a9bb820960c5b93c9fdec209c0360a&",
	"https://cdn.discordapp.com/attachments/1112127464584511538/1261031419619184757/12.mp4?ex=669374c7&is=66922347&hm=7af1d947e6550a63675ca97e0dd2177b04268c7095b0b041dc05ada6f1a13742&",
	"https://cdn.discordapp.com/attachments/1112127464584511538/1261031420369833994/15.mov?ex=669374c7&is=66922347&hm=98004e07cf03701a5338f9c5a267283a6ade641a2dafabceb06a8c1b05ed70a3&",
	"https://cdn.discordapp.com/attachments/1112127464584511538/1261031415311499316/SPOILER_videoplayback.mp4?ex=669374c6&is=66922346&hm=23a4c73c75219de5932569ef60a17ea752d298eb5dcd8de3997d9d5041969375&",
	"https://cdn.discordapp.com/attachments/1112127464584511538/1261031416368730285/statham_mean_machine_2.mov?ex=669374c6&is=66922346&hm=01d482efc41cdfc0b49454cc718c8bcd3456ee7ab13c596cef2141b01e854639&",
	"https://cdn.discordapp.com/attachments/1112127464584511538/1261031415869472929/statham_mean_machiene_1.mov?ex=669374c6&is=66922346&hm=91e6f0fa3b05f5a1aa981ad109783eb0d39e355183fee55a18a302b67570f382&",
	"https://cdn.discordapp.com/attachments/1112127464584511538/1261031393308315713/11.mp4?ex=669374c0&is=66922340&hm=fabfa1174c74502017a944a900e0a85796118dccd82ce15cbcab12790150a666&",
	"https://cdn.discordapp.com/attachments/1112127464584511538/1261031392196956160/7bdd66b4f3e5a477.mp4?ex=669374c0&is=66922340&hm=d59ade7a74164fe8d231545b70196bb02b8c7bf11c0bb9beacf63f96e2e16d39&",
	"https://cdn.discordapp.com/attachments/1112127464584511538/1261031392763052032/9.mov?ex=669374c0&is=66922340&hm=4dc633b7cd1f1bdb33911cb875af271770284183a5eac73972968232e907a05c&",
	"https://cdn.discordapp.com/attachments/1112127464584511538/1261031356406829076/7.mp4?ex=669374b8&is=66922338&hm=526aac9db1d747d6cdc1ce966c56130d31f374936720493f1cbb5817521e3794&",
	"https://cdn.discordapp.com/attachments/1112127464584511538/1261031355463106571/4.mp4?ex=669374b7&is=66922337&hm=905ea122207d92141a4981be11930edaa98976525cfacf3d8a0427105dfbddaa&",
	"https://cdn.discordapp.com/attachments/1112127464584511538/1261031355895251128/6.mov?ex=669374b8&is=66922338&hm=ff51fbac91256ae1b76ce8a277b12817e77a5a2e3ebdc129bc085e546d2afb3d&",
	"https://cdn.discordapp.com/attachments/1112127464584511538/1261031295085969500/4.mov?ex=669374a9&is=66922329&hm=1680961e569671f3409fd8686e0b5ffc79c6a4ac750fc0fd6767d45ef83cfb4a&",
	"https://cdn.discordapp.com/attachments/1112127464584511538/1261031293974479010/1.mov?ex=669374a9&is=66922329&hm=e85352117f243ed87c0053c673742c5bb6b24b457f9b9e4a698b09b87383c5ec&",
	"https://cdn.discordapp.com/attachments/1112127464584511538/1261031294624731136/2.mp4?ex=669374a9&is=66922329&hm=c589e1cace33c464913a39e5add4424c54dda6e24391828089fb7705daba5520&",
}

func RequestJason(message *discordgo.MessageCreate) {
	SendJason(message, "")
	dbhelper.CountCommand("requestjason", message.Author.ID)
}

func DetectJason(message *discordgo.MessageCreate) bool {
	if strings.Contains(strings.ToLower(message.Content), "jason") {
		SendJason(message, "Jason")
		dbhelper.CountCommand("detectjason", message.Author.ID)
		return true
	} else if strings.Contains(strings.ToLower(message.Content), "statham") {
		SendJason(message, "Statham")
		dbhelper.CountCommand("detectjason", message.Author.ID)
		return true
	} else {
		return false
	}
}

func SendJason(message *discordgo.MessageCreate, foundInText string) {

	if foundInText == "" {
		foundInText = "."
	} else {
		foundInText = foundInText + "?"
	}
	randomJason := helpers.GetRandomText(jasonVideos)
	_, err := config.Session.ChannelMessageSend(message.ChannelID, "["+foundInText+"]("+randomJason+")")
	if err != nil {
		audit.Error(err)
		logging.SendError(message)
	}

}
