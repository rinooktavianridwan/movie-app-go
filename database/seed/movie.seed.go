package seed

import (
	"movie-app-go/internal/models"

	"gorm.io/gorm"
)

func SeedMovies(db *gorm.DB) ([]models.Movie, error) {
	movies := []models.Movie{
		{Title: "Avengers: Endgame", Overview: "After the devastating events of Avengers: Infinity War, the universe is in ruins due to the efforts of the Mad Titan, Thanos. With the help of remaining allies, the Avengers must assemble once more in order to undo Thanos's actions and restore order to the universe once and for all.", Duration: 181},
		{Title: "Spider-Man: No Way Home", Overview: "Peter Parker is unmasked and no longer able to separate his normal life from the high-stakes of being a super-hero. When he asks for help from Doctor Strange the stakes become even more dangerous, forcing him to discover what it truly means to be Spider-Man.", Duration: 148},
		{Title: "The Dark Knight", Overview: "Batman raises the stakes in his war on crime. With the help of Lt. Jim Gordon and District Attorney Harvey Dent, Batman sets out to dismantle the remaining criminal organizations that plague the streets.", Duration: 152},
		{Title: "Inception", Overview: "Dom Cobb is a skilled thief, the absolute best in the dangerous art of extraction, stealing valuable secrets from deep within the subconscious during the dream state, when the mind is at its most vulnerable.", Duration: 148},
		{Title: "The Shawshank Redemption", Overview: "Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.", Duration: 142},
		{Title: "Interstellar", Overview: "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.", Duration: 169},
		{Title: "The Godfather", Overview: "The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.", Duration: 175},
		{Title: "Pulp Fiction", Overview: "The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.", Duration: 154},
		{Title: "Forrest Gump", Overview: "The presidencies of Kennedy and Johnson, the Vietnam War, the Watergate scandal and other historical events unfold from the perspective of an Alabama man with an IQ of 75.", Duration: 142},
		{Title: "The Matrix", Overview: "A computer hacker learns from mysterious rebels about the true nature of his reality and his role in the war against its controllers.", Duration: 136},
		{Title: "Goodfellas", Overview: "The story of Henry Hill and his life in the mob, covering his relationship with his wife Karen Hill and his mob partners Jimmy Conway and Tommy DeVito.", Duration: 146},
		{Title: "The Lord of the Rings: The Fellowship of the Ring", Overview: "A meek Hobbit from the Shire and eight companions set out on a journey to destroy the powerful One Ring and save Middle-earth from the Dark Lord Sauron.", Duration: 178},
		{Title: "Star Wars: A New Hope", Overview: "Luke Skywalker joins forces with a Jedi Knight, a cocky pilot, a Wookiee and two droids to save the galaxy from the Empire's world-destroying battle station.", Duration: 121},
		{Title: "Fight Club", Overview: "An insomniac office worker and a devil-may-care soapmaker form an underground fight club that evolves into something much, much more.", Duration: 139},
		{Title: "The Lion King", Overview: "A young lion prince is cast out of his pride by his cruel uncle, who claims he killed his father so that he can become the new king.", Duration: 88},
		{Title: "Toy Story", Overview: "A cowboy doll is profoundly threatened and jealous when a new spaceman figure supplants him as top toy in a boy's room.", Duration: 81},
		{Title: "Jurassic Park", Overview: "A pragmatic paleontologist visiting an almost complete theme park is tasked with protecting a couple of kids after a power failure causes the park's cloned dinosaurs to run loose.", Duration: 127},
		{Title: "Titanic", Overview: "A seventeen-year-old aristocrat falls in love with a kind but poor artist aboard the luxurious, ill-fated R.M.S. Titanic.", Duration: 194},
		{Title: "The Silence of the Lambs", Overview: "A young F.B.I. cadet must receive the help of an incarcerated and manipulative cannibal killer to help catch another serial killer, a madman who skins his victims.", Duration: 118},
		{Title: "Saving Private Ryan", Overview: "Following the Normandy Landings, a group of U.S. soldiers go behind enemy lines to retrieve a paratrooper whose brothers have been killed in action.", Duration: 169},
		{Title: "Schindler's List", Overview: "In German-occupied Poland during World War II, industrialist Oskar Schindler gradually becomes concerned for his Jewish workforce after witnessing their persecution by the Nazis.", Duration: 195},
		{Title: "La La Land", Overview: "While navigating their careers in Los Angeles, a pianist and an actress fall in love while attempting to reconcile their aspirations for the future.", Duration: 128},
		{Title: "Parasite", Overview: "Greed and class discrimination threaten the newly formed symbiotic relationship between the wealthy Park family and the destitute Kim clan.", Duration: 132},
		{Title: "Joker", Overview: "In Gotham City, mentally troubled comedian Arthur Fleck is disregarded and mistreated by society. He then embarks on a downward spiral of revolution and bloody crime.", Duration: 122},
		{Title: "Black Panther", Overview: "T'Challa, heir to the hidden but advanced kingdom of Wakanda, must step forward to lead his people into a new future and must confront a challenger from his country's past.", Duration: 134},
		{Title: "Frozen", Overview: "When the newly crowned Queen Elsa accidentally uses her power to turn things into ice to curse her home in infinite winter, her sister Anna teams up with a mountain man, his playful reindeer, and a snowman to change the weather condition.", Duration: 102},
		{Title: "Finding Nemo", Overview: "After his son is captured in the Great Barrier Reef and taken to Sydney, a timid clownfish sets out on a journey to bring him home.", Duration: 100},
		{Title: "The Incredibles", Overview: "A family of undercover superheroes, while trying to live the quiet suburban life, are forced into action to save the world.", Duration: 115},
		{Title: "WALL-E", Overview: "In the distant future, a small waste-collecting robot inadvertently embarks on a space journey that will ultimately decide the fate of mankind.", Duration: 98},
		{Title: "Up", Overview: "78-year-old Carl Fredricksen travels to Paradise Falls in his house equipped with balloons, inadvertently taking a young stowaway.", Duration: 96},
		{Title: "Inside Out", Overview: "After young Riley is uprooted from her Midwest life and moved to San Francisco, her emotions - Joy, Fear, Anger, Disgust and Sadness - conflict on how best to navigate a new city, house, and school.", Duration: 95},
	}

	if err := db.Create(&movies).Error; err != nil {
		return nil, err
	}
	return movies, nil
}
