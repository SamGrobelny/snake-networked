import pygame
from enum import Enum
import random
import numpy as np

class Direction(Enum):
    NONE = 0
    UP = 1
    DOWN = 2
    LEFT = 3
    RIGHT = 4
    
GRID_WIDTH = 20
GRID_HEIGHT = 20

grid = np.zeros((GRID_WIDTH, GRID_HEIGHT), int)

screen_width = 1280
screen_height = 640

cell_width = screen_width // GRID_WIDTH
cell_height = screen_height // GRID_HEIGHT
cell_size = min(cell_height, cell_width)

pygame.init()

screen = pygame.display.set_mode((screen_width, screen_height))
running = True
clock = pygame.time.Clock()

# player default properties
player_direction = Direction.NONE
player_x = GRID_WIDTH / 2
player_y = GRID_HEIGHT / 2
player_pos = pygame.Vector2(player_x, player_y)

snake_arr = [player_pos]

player_speed = 10.0 # in cells per second
delta_req = 1 / player_speed
time_accum = 0


# food
def respawn_food():
    
    food_x = random.randint(0,GRID_WIDTH-1)
    food_y = random.randint(0, GRID_HEIGHT-1)
    pos = pygame.Vector2(food_x, food_y)
    
    while pos == snake_arr[-1]:
        food_x = random.randint(0,GRID_WIDTH-1)
        food_y = random.randint(0, GRID_HEIGHT-1)
        pos = pygame.Vector2(food_x, food_y)
        
    return pos

food_pos = respawn_food()


def reset_player():
    pos = pygame.Vector2(player_x, player_y)
    d = Direction.NONE
    
    return pos, d

GAME_FONT = pygame.freetype.Font("Ldfcomicsans-jj7l.ttf", 16)

# game loop
while running:
    time_delta = clock.tick(60) / 1000.0 # convert to seconds
    
    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            running = False
            
    # change direction based on key pressed
    keys = pygame.key.get_pressed()
    
    if keys[pygame.K_UP]:
        player_direction = Direction.UP
    if keys[pygame.K_DOWN]:
        player_direction = Direction.DOWN
    if keys[pygame.K_LEFT]:
        player_direction = Direction.LEFT
    if keys[pygame.K_RIGHT]:
        player_direction = Direction.RIGHT
    
    # calculate if snake moves
    time_accum += time_delta
    if time_accum >= delta_req:
        # move snake here
        match player_direction:
            case Direction.UP:
                snake_arr[-1].y -= 1
            case Direction.DOWN:
                snake_arr[-1].y += 1
            case Direction.LEFT:
                snake_arr[-1].x -= 1
            case Direction.RIGHT:
                snake_arr[-1].x += 1
                
        if player_direction != Direction.NONE and len(snake_arr) > 1:
            for i in range(len(snake_arr) - 1):
                snake_arr[i] = snake_arr[i+1]
        
        time_accum -= delta_req
        
    # calculate if out of bound
    # reset player if so
    if snake_arr[-1].x >= GRID_WIDTH or snake_arr[-1].x < 0:
        pos, player_direction = reset_player()
        snake_arr = [pos]
    elif snake_arr[-1].y >= GRID_HEIGHT or snake_arr[-1].y < 0:
        pos, player_direction = reset_player()
        snake_arr = [pos]
        
    if snake_arr[-1] == food_pos:
        snake_arr.append(food_pos)
        food_pos = respawn_food()
        
    # drawing time!
    screen.fill("black")
    
    # draw background
    bg = pygame.Rect(
        0,
        0,
        (GRID_WIDTH) * cell_size,
        (GRID_HEIGHT) * cell_size
    )
    
    # draw snake here
    # snake = pygame.Rect(
    #     player_pos.x * cell_size,
    #     player_pos.y * cell_size,
    #     cell_size,
    #     cell_size
    # )
    
    def draw_snake_cell(pos : pygame.Vector2) -> pygame.Rect:
        return pygame.Rect(
            pos.x * cell_size,
            pos.y * cell_size,
            cell_size,
            cell_size
        )
        
    
    # draw food here
    food = pygame.Rect(
        food_pos.x * cell_size,
        food_pos.y * cell_size,
        cell_size,
        cell_size
    )
    
    pygame.draw.rect(screen, "grey10", bg)
    
    pygame.draw.rect(screen, "green", food)
    
    for cell in snake_arr:
        cell_rect = draw_snake_cell(cell)
        pygame.draw.rect(screen, "red", cell_rect)
    
    GAME_FONT.render_to(
        screen,
        (2, 2),
        f'arr:{snake_arr}',
        (255, 255, 255)
    )
    
    pygame.display.flip()