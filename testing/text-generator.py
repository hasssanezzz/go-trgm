from random import choice, choices, randint
from sys import argv

# Keep the original ALPHA
ALPHA = "abcdefghijklmnopqrstuvwxyz0123456789 !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"

# Define letter frequencies based on English text (approximation)
# Source: https://en.wikipedia.org/wiki/Letter_frequency
# Weights are roughly proportional to frequency percentages.
# Spaces and punctuation are added with lower weights.
LETTER_FREQ = {
    'e': 12.7, 't': 9.1, 'a': 8.2, 'o': 7.5, 'i': 7.0, 'n': 6.7, 's': 6.3,
    'h': 6.1, 'r': 6.0, 'd': 4.3, 'l': 4.0, 'c': 2.8, 'u': 2.8, 'm': 2.4,
    'w': 2.4, 'f': 2.2, 'g': 2.0, 'y': 2.0, 'p': 1.9, 'b': 1.3, 'v': 1.0,
    'k': 0.8, 'j': 0.15, 'x': 0.15, 'q': 0.10, 'z': 0.07,
    ' ': 15.0, # High frequency for spaces
    # Add some common punctuation with low frequency
    '.': 1.0, ',': 1.0, '!': 0.5, '?': 0.5, "'": 0.5, '"': 0.2
    # Numbers and other symbols from ALPHA can be included with very low/default weight
}

# Create a list of letters and corresponding weights for choices()
# Include all letters from ALPHA, assigning default low weight if not in LETTER_FREQ
default_weight = 0.05
letters = list(ALPHA)
weights = [LETTER_FREQ.get(char, default_weight) for char in letters]

# Function to pick a letter based on frequency
def letter():
    return choices(letters, weights=weights, k=1)[0]

# --- Word Generation ---
# Adjusted to make more common word lengths slightly more probable
# and to make short/long word distinction less rigid.

# Generate a pool of 'words' with varied lengths
# Mostly short words (1-5 chars), some medium (6-10), fewer long (10+)
word_pool = []
for _ in range(500): # Increased pool size slightly
    # Skew towards shorter lengths
    if randint(1, 10) <= 6:  # 60% chance for short words (1-5)
        length = randint(1, 5)
    elif randint(1, 10) <= 9: # 30% chance for medium words (6-10)
        length = randint(6, 10)
    else:                     # 10% chance for longer words (11-15)
        length = randint(11, 15)
    word_pool.append(''.join(letter() for _ in range(length)))

# --- Sentence Generation ---
def sentence():
    # Vary sentence length
    num_words = randint(5, 25) # Slightly wider range
    words = []
    for i in range(num_words):
        words.append(choice(word_pool))
        # Add a chance for basic punctuation (comma, period)
        # Mostly commas in the middle, period at the end
        if i < num_words - 1 and randint(1, 20) == 1: # 5% chance for comma before last word
            # Insert comma logically (before the next word)
            words[-1] += ',' # Add comma to the last added word
        elif i == num_words - 1: # Last word
             # Add period, exclamation, or question mark
             ending = choice(['.', '.', '.', '!', '?']) # Period more likely
             words[-1] += ending

    sentence_text = ' '.join(words)
    # Basic capitalization
    if sentence_text:
        sentence_text = sentence_text[0].upper() + sentence_text[1:]
    return sentence_text

# --- Main Execution ---
def main():
    if len(argv) != 2:
        print("Usage: python script.py <number_of_sentences>")
        return

    try:
        num_sentences = int(argv[1])
    except ValueError:
        print("Error: The argument must be an integer.")
        return

    # Generate sentences
    sentences = [sentence() for _ in range(num_sentences)]

    # Write to file
    try:
        with open("lines.better.txt", 'w') as f:
            f.write('\n'.join(sentences))
        print(f"Generated {num_sentences} sentences and wrote to 'lines.txt'")
    except IOError as e:
        print(f"Error writing to file: {e}")

if __name__ == "__main__":
    main()
