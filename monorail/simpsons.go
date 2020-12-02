package main

// Simpsons character array
// We'll use this as a quick way to get the profile drawing for each character
var simpsonsCharacters = map[string]SimpsonCharacter{
	"homer": {
		Name:      "homer",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/0/02/Homer_Simpson_2006.png",
		Quote:     "Mmm, donuts.",
	},
	"marge": {
		Name:      "marge",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/0/0b/Marge_Simpson.png",
		Quote:     "I don't think that's a very good idea.",
	},
	"bart": {
		Name:      "bart",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/a/aa/Bart_Simpson_200px.png",
		Quote:     "I'm Bart Simpson, Who the Hell are You?",
	},
	"lisa": {
		Name:      "lisa",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/e/ec/Lisa_Simpson.png",
		Quote:     "If anyone wants me, I'll be in my room.",
	},
	"maggie": {
		Name:      "maggie",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/9/9d/Maggie_Simpson.png",
		Quote:     "(takes two steps, immediately falls)",
	},
	"abe": {
		Name:      "abe",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/3/3e/Abe_Simpson.png",
		Quote:     "Hot Diggity Dog!",
	},
	"apu": {
		Name:      "apu",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/2/23/Apu_Nahasapeemapetilon_%28The_Simpsons%29.png",
		Quote:     "No sir. In fact, I can recite pi to forty thousand places. The last digit is 1.",
	},
	"barney": {
		Name:      "barney",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/d/de/Barney_Gumble.png",
		Quote:     "(Burp!)",
	},
	"wiggum": {
		Name:      "wiggum",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/7/7a/Chief_Wiggum.png",
		Quote:     "Cuff 'em, Lou",
	},
	"lou": {
		Name:      "lou",
		AvatarURL: "https://static.wikia.nocookie.net/simpsons/images/1/17/Lou.png",
		Quote:     "Uh... Chief?",
	},
	"krusty": {
		Name:      "krusty",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/5/5a/Krustytheclown.png",
		Quote:     "Hey, Hey Kids!",
	},
	"milhouse": {
		Name:      "milhouse",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/1/11/Milhouse_Van_Houten.png",
		Quote:     "Everything's coming up Milhouse!",
	},
	"moe": {
		Name:      "moe",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/8/80/Moe_Szyslak.png",
		Quote:     "Moe's Tavern, Moe speaking.",
	},
	"burns": {
		Name:      "mr. burns",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/5/56/Mr_Burns.png",
		Quote:     "Excellent!",
	},
	"ned": {
		Name:      "ned",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/8/84/Ned_Flanders.png",
		Quote:     "Hi-Diddily-Ho Neighborino!",
	},
	"ralph": {
		Name:      "ralph",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/1/14/Ralph_Wiggum.png",
		Quote:     "Haha! I'm in danger!",
	},
	"lovejoy": {
		Name:      "lovejoy",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/7/7d/Rev_Lovejoy.png",
		Quote:     "Wait a minute, this sounds like rock and or roll!",
	},
	"patty": {
		Name:      "patty",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/f/f8/Patty_Bouvier.png",
		Quote:     "Some days at the DMV, we don't let the line move at all. We call those weekdays.",
	},
	"selma": {
		Name:      "selma",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/b/ba/Selma_Bouvier.png",
		Quote:     "How's my blubber in-law?",
	},
	"skinner": {
		Name:      "skinner",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/3/3a/Seymour_Skinner.png",
		Quote:     "I know very little about children.",
	},
	"bob": {
		Name:      "bob",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/c/c8/C-bob.png",
		Quote:     "Hello, Bart.",
	},
	"smithers": {
		Name:      "smithers",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/8/86/Waylon_Smithers_1.png",
		Quote:     "That's Homer Simpson, one of your blobs from sector 7G.",
	},
	"edna": {
		Name:      "edna",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/7/76/Edna_Krabappel.png",
		Quote:     "Ha!",
	},
	"quimby": {
		Name:      "quimby",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/5/51/Mayor_Quimby.png",
		Quote:     "Vote Quimby!",
	},
	"nelson": {
		Name:      "nelson",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/c/c6/Nelson_Muntz.PNG",
		Quote:     "Ha Ha!",
	},
	"frink": {
		Name:      "frink",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/7/71/Frink.png",
		Quote:     "Oops, I forgot to carry the one.",
	},
	"troy": {
		Name:      "troy mcclure",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/6/6c/Troymcclure.png",
		Quote:     "Hi, I'm Troy McClure!",
	},
	"comic_book_guy": {
		Name:      "comic book guy",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/7/79/The_Simpsons-Jeff_Albertson.png",
		Quote:     "Worst Website Ever!",
	},
	"willie": {
		Name:      "willie",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/d/dc/GroundskeeperWillie.png",
		Quote:     "Willie hears ya, Willie don't care.",
	},
	"fat_tony": {
		Name:      "fat tony",
		AvatarURL: "https://upload.wikimedia.org/wikipedia/en/3/3e/FatTony.png",
		Quote:     "Uh, I Prefer The Cat. He Hates Mondays. We Can All Relate.",
	},
	"hibbert": {
		Name:      "dr. hibbert",
		AvatarURL: "https://static.wikia.nocookie.net/simpsons/images/6/67/Tapped_Out_Unlock_DrHibbert.png",
		Quote:     "Have a wowwipop.",
	},
}

// The array is used to create the mapping between your profile image on github, and the character
var characterIndex = []string{
	"ralph", "homer", "marge", "bart", "lisa",
	"maggie", "abe", "apu", "barney", "krusty",
	"milhouse", "moe", "burns", "ned", "lovejoy",
	"patty", "selma", "skinner", "wiggum", "lou",
	"bob", "smithers", "edna", "quimby", "nelson",
	"frink", "troy", "comic_book_guy", "willie",
	"fat_tony", "hibbert",
}
