import pygame
from enum import Enum
import random
import numpy as np


class SnakeNode:
    def __init__(self, value, next=None):
        self.value = value
        self.next = next


class LinkedList:
    def __init__(self, head: SnakeNode):
        self.head = head
        self.tail = None

    def get_head(self):
        return self.head

    def set_head(self, value):
        temp = SnakeNode(value=value)
        temp.next = self.head
        self.head = temp

    def remove_tail(self):
        curr = self.head
        while curr.next is not self.tail:
            curr = curr.next
        curr.next = None
        self.tail = curr


class Snake:
    # TODO: create the very pogger snake class


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

offset = (screen_width - (GRID_WIDTH * cell_size)) // 2

pygame.init()

screen = pygame.display.set_mode((screen_width, screen_height))
running = True
clock = pygame.time.Clock()

# player default properties
player_direction = Direction.NONE
player_x = GRID_WIDTH / 2
player_y = GRID_HEIGHT / 2
player_pos = pygame.Vector2(player_x, player_y)

snake_list = LinkedList(head=SnakeNode(value=player_pos))

player_speed = 9.0  # in cells per second
delta_req = 1 / player_speed
time_accum = 0


def increment_snake(direction, head: LinkedList):
    x = head.value.x
    y = head.value.y

    match direction:
        case Direction.UP:
            # player_pos.y -= 1
            y -= 1
        case Direction.DOWN:
            # player_pos.y += 1
            y += 1
        case Direction.LEFT:
            # player_pos.x -= 1
            x -= 1
        case Direction.RIGHT:
            # player_pos.x += 1
            x += 1

    new_head = pygame.Vector2(x, y)
    head.set_head(new_head)
    head.remove_tail()

# food


def respawn_food():

    food_x = random.randint(0, GRID_WIDTH-1)
    food_y = random.randint(0, GRID_HEIGHT-1)
    pos = pygame.Vector2(food_x, food_y)

    while pos == player_pos:
        food_x = random.randint(0, GRID_WIDTH-1)
        food_y = random.randint(0, GRID_HEIGHT-1)
        pos = pygame.Vector2(food_x, food_y)

    return pos


food_pos = respawn_food()


def reset_player():
    pos = pygame.Vector2(player_x, player_y)
    d = Direction.NONE

    return pos, d


# game loop
while running:
    time_delta = clock.tick(60) / 1000.0  # convert to seconds

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
        increment_snake(player_direction, snake_head)
        # TODO: replace with increment player function
        match player_direction:
            case Direction.UP:
                player_pos.y -= 1
            case Direction.DOWN:
                player_pos.y += 1
            case Direction.LEFT:
                player_pos.x -= 1
            case Direction.RIGHT:
                player_pos.x += 1

        time_accum -= delta_req

    # calculate if out of bound
    # reset player if so
    if player_pos.x >= GRID_WIDTH or player_pos.x < 0:
        player_pos, player_direction = reset_player()
    elif player_pos.y >= GRID_HEIGHT or player_pos.y < 0:
        player_pos, player_direction = reset_player()

    if player_pos == food_pos:
        food_pos = respawn_food()

    # drawing time!
    screen.fill("black")

    # draw background
    bg = pygame.Rect(
        offset,
        0,
        (GRID_WIDTH) * cell_size,
        (GRID_HEIGHT) * cell_size
    )

    # draw food here
    food = pygame.Rect(
        offset + (food_pos.x * cell_size),
        food_pos.y * cell_size,
        cell_size,
        cell_size
    )

    pygame.draw.rect(screen, "grey10", bg)
    pygame.draw.rect(screen, "green", food)

    # draw snake here
    curr = snake_head
    while curr.next is not None:
        rect = pygame.Rect(
            offset + (curr.value.x * cell_size),
            curr.value.y * cell_size,
            cell_size,
            cell_size
        )
        pygame.draw.rect(screen, "red", rect)

    pygame.display.flip()
