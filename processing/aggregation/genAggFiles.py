"""
Generates all of the files to do with running an N dimensional aggregate, except for aggND Should be run before the code
is compiled as produced methods are much faster than aggND. Functions similar to aggND but with loops unwrapped and the
array of coordinates stored as variables in the struct
"""

import os

MAX_DIMENSION = 7  # can be set to any positive value, files will be generated up to and including this value

aggND = '//This is an auto generated file from genAggFiles.py\n' \
        'package agg{n}D\n' \
        '\n' \
        'import (\n' \
        '\t"encoding/gob"\n' \
        '\t"math"\n' \
        '\t"math/rand"\n' \
        ')\n' \
        '\n' \
        'func init() {{\n' \
        '\tgob.Register(Point{n}D{{}})\n' \
        '}}\n' \
        '\n' \
        'const (\n' \
        '\tBORDER_SCALE float64 = {border_scale}\n' \
        '\tBORDER_CONST float64 = {border_const}\n' \
        ')\n' \
        '\n' \
        'type Point{n}D struct {{\n' \
        '\t{letters} int64\n' \
        '}}\n' \
        '\n' \
        'func (p Point{n}D) Coordinates() []int64 {{\n' \
        '\treturn []int64{{ {p_letters} }}\n' \
        '}}\n' \
        '\n' \
        'func (p Point{n}D) SquareDistance(coords []float64) float64 {{\n' \
        '\tvar {d_letters} = {d_coords}\n' \
        '\treturn {sq_dist}\n' \
        '}}\n' \
        '\n' \
        '\n' \
        'type cache struct {{\n' \
        '\tpoint Point{n}D\n' \
        '\tpointRadius float64\n' \
        '\trng *rand.Rand\n' \
        '\tlastWalk int64\n' \
        '\tstate map[Point{n}D]int64\n' \
        '\tstateRadius float64\n' \
        '\tborderRadius float64\n' \
        '\tborderRadiusInt int64\n' \
        '\ttempPoint Point{n}D\n' \
        '\ttempFloatA, tempFloatB float64\n' \
        '}}\n' \
        '\n' \
        'func (c *cache) updateCurrPointRadius() {{\n' \
        '\tc.pointRadius = math.Sqrt(float64({c_sq_dist}))\n' \
        '}}\n' \
        '\n' \
        'func (c *cache) updateStateRadius() {{\n' \
        '\tc.stateRadius = c.pointRadius\n' \
        '\tc.borderRadius = c.stateRadius*BORDER_SCALE+BORDER_CONST\n' \
        '\tc.borderRadiusInt = int64(c.borderRadius)\n' \
        '}}\n' \
        '\n' \
        'func (c *cache) pointIn() (ok bool) {{\n' \
        '\t_, ok = c.state[c.point]\n' \
        '\treturn\n' \
        '}}\n' \
        '\n' \
        '\n' \
        'func (c *cache) pointToBorder() {{\n' \
        '{mov_border}\n' \
        '\tc.updateCurrPointRadius()\n' \
        '}}\n' \
        '\n' \
        'func (c *cache) walkPoint() {{\n' \
        '\t{walk_switch}\n' \
        '\tc.updateCurrPointRadius()\n' \
        '}}\n' \
        '\n' \
        '{check_neighbors}' \
        '\n' \
        'func RunNew(nPoints int64, sticking float64, rng *rand.Rand) map[Point{n}D]int64 {{\n' \
        '\tc := cache{{}}\n' \
        '\tc.rng = rng\n' \
        '\tc.state = make(map[Point{n}D]int64, nPoints)\n' \
        '\tc.state[c.point] = 0\n' \
        '\tc.updateStateRadius()\n' \
        '\t\n' \
        '\tfor i := int64(1); i < nPoints; i++ {{\n' \
        '\t\tc.pointToBorder()\n' \
        '\t\tfor !c.pointHasNeighbor() || sticking < rng.Float64() {{\n' \
        '\t\t\tc.walkPoint()\n' \
        '\t\t}}\n' \
        '\t\tc.state[c.point]=i\n' \
        '\t\tif c.pointRadius > c.stateRadius {{\n' \
        '\t\t\tc.updateStateRadius()\n' \
        '\t\t}}\n' \
        '\t}}\n' \
        '\treturn c.state\n' \
        '}}\n'
aggregaton = '//This is an auto generated file from genAggFiles.py\n' \
             'package aggregation\n' \
             '\n' \
             'import (\n' \
             '\t"math/rand"\n' \
             '{imports}\n' \
             '\t"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/aggND"\n' \
             ')\n' \
             '\n' \
             'type Point interface {{\n' \
             '\tCoordinates() []int64\n' \
             '\tSquareDistance([]float64) float64\n' \
             '}}\n' \
             '\n' \
             '{runs}\n' \
             'func RunND(nPoints, seed int64, sticking float64, nDimension int64) (points []Point) {{\n' \
             '\tvar rng = rand.New(rand.NewSource(seed))\n' \
             '\tvar state = aggND.RunNew(nPoints, sticking, rng, nDimension)\n' \
             '\tif int64(len(state)) != nPoints {{\n' \
             '\t\tpanic("N != nPoints. This should never happen.")\n' \
             '\t}}\n' \
             '\tpoints = make([]Point, nPoints)\n' \
             '\tfor i, p := range state {{\n' \
             '\t\tpoints[i] = p\n' \
             '\t}}\n' \
             '\treturn points\n' \
             '}}' \
             '\n' \
             'func RunNew(nPoints, seed, nDimension int64, sticking float64) []Point {{\n' \
             '\tswitch nDimension {{\n' \
             '{cases}\n' \
             '\tdefault: return RunND(nPoints, seed, sticking, nDimension)\n' \
             '\t}}\n' \
             '}}\n'
run = 'func Run{n}D(nPoints, seed int64, sticking float64) (points []Point) {{\n' \
      '\tvar rng = rand.New(rand.NewSource(seed))\n' \
      '\tvar state = agg{n}D.RunNew(nPoints, sticking, rng)\n' \
      '\tif int64(len(state)) != nPoints {{\n' \
      '\t\tpanic("N != nPoints. This should never happen.")\n' \
      '\t}}\n' \
      '\tpoints = make([]Point, nPoints)\n' \
      '\tfor point, index := range state {{\n' \
      '\t\tpoints[index] = point\n' \
      '\t}}\n' \
      '\treturn points\n' \
      '}}\n'

def get_next_coord(s):
    if len(s) == 0:
        return 'A'
    e = s[-1]
    if e == 'Z':
        return get_next_coord(s[:-1]) + 'A'
    else:
        e = chr(ord(e)+1)
    return s[:-1] + e

def get_letters(n):
    x = ''
    letters = []
    for i in range(n):
        x = get_next_coord(x)
        letters += x
    return letters

def get_mov_border(letters):
    ret = '\tc.tempFloatA = 1\n'
    for x in letters[:-1]:
        ret += '\tc.tempFloatB = c.rng.Float64() * 2 * math.Pi\n'
        ret += '\tc.point.{} = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.borderRadius)\n'.format(x)
        ret += '\tc.tempFloatA *= math.Sin(c.tempFloatB)\n'
    ret += '\tc.point.{} = int64(c.tempFloatA * c.borderRadius)\n'.format(letters[-1])
    return ret

def get_walk_switch(letters):
    ret = '\tswitch c.rng.Int63n({}) {{\n'.format(len(letters)*2)
    for i, l in enumerate(letters):
        ret += '\tcase {na}:\n' \
               '\t\tc.point.{l}++\n' \
               '\t\tif c.pointRadius < c.stateRadius+BORDER_CONST && c.pointIn() {{\n' \
               '\t\t\tc.point.{l}--\n' \
               '\t\t}} else {{\n' \
               '\t\t\tif c.point.{l} > c.borderRadiusInt {{\n' \
               '\t\t\t\tc.point.{l} -= 2*c.borderRadiusInt\n' \
               '\t\t\t}}\n' \
               '\t\t\tc.lastWalk = {na}\n' \
               '\t\t}}\n' \
               '\tcase {nb}:\n' \
               '\t\tc.point.{l}--\n' \
               '\t\tif c.pointRadius < c.stateRadius+BORDER_CONST && c.pointIn() {{\n' \
               '\t\t\tc.point.{l}++\n' \
               '\t\t}} else {{\n' \
               '\t\t\tif c.point.{l} < -c.borderRadiusInt {{\n' \
               '\t\t\t\tc.point.{l} += 2*c.borderRadiusInt\n' \
               '\t\t\t}}\n' \
               '\t\t\tc.lastWalk = {nb}\n' \
               '\t\t}}\n'.format(na=i*2, nb=i*2 + 1, l=l)
    return ret + "}\n"

def get_check_neighbors(letters):
    ret = ''
    check = []
    for i, l in enumerate(letters):
        ret += 'func (c *cache) isPlus{l}In() bool {{\n' \
               '\treturn c.lastWalk != {nb} && c.plus{l}In()\n' \
               '}}\n' \
               'func (c *cache) plus{l}In() (ok bool) {{\n' \
               '\tc.tempPoint = c.point\n' \
               '\tc.tempPoint.{l}++\n' \
               '\t_, ok = c.state[c.tempPoint]\n' \
               '\treturn\n' \
               '}}\n' \
               'func (c *cache) isMinus{l}In() bool {{\n' \
               '\treturn c.lastWalk != {na} && c.minus{l}In()\n' \
               '}}\n' \
               'func (c *cache) minus{l}In() (ok bool) {{\n' \
               '\tc.tempPoint = c.point\n' \
               '\tc.tempPoint.{l}--\n' \
               '\t_, ok = c.state[c.tempPoint]\n' \
               '\treturn\n' \
               '}}\n\n'.format(l=l, na=i*2, nb=i*2 + 1)
        check += ['c.isPlus{l}In() || c.isMinus{l}In()'.format(l=l)]
    ret += 'func (c *cache) pointHasNeighbor() bool {{\n' \
           '\treturn c.pointRadius <= c.stateRadius+1 && ({})\n' \
           '}}\n'.format(' || '.join(check))
    return ret

def make(n: int):
    letters = get_letters(n)
    fields = {
        'n': n,
        'border_scale': 1.5,
        'border_const': 2,
        'letters': ', '.join(letters),
        'p_letters': ', '.join(map(lambda x: 'p.{}'.format(x), letters)),
        'd_letters': ', '.join(map(lambda x: 'd{}'.format(x), letters)),
        'd_coords': ', '.join(map(lambda x: 'float64(p.{1})-coords[{0}]'.format(*x), enumerate(letters))),
        'sq_dist': ' + '.join(map(lambda x: 'd{0}*d{0}'.format(x), letters)),
        'c_sq_dist': ' + '.join(map(lambda x: 'c.point.{0}*c.point.{0}'.format(x), letters)),
        'mov_border': get_mov_border(letters),
        'walk_switch': get_walk_switch(letters),
        'check_neighbors': get_check_neighbors(letters),
    }
    return aggND.format(**fields)

def make_aggregation(entries):
    p = '\t"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg{n}D"'
    c = '\tcase {n}: return Run{n}D(nPoints, seed, sticking)\n'
    fields = {
        "imports": '\n'.join(map(lambda x: p.format(n=x), entries)),
        "runs": '\n'.join(map(lambda x: run.format(n=x), entries)),
        "cases": ''.join(map(lambda x: c.format(n=x), entries)),
    }
    return aggregaton.format(**fields)

def main():
    entries = []
    for i in range(1, MAX_DIMENSION+1):
        if not os.path.exists("agg{0}D".format(i)):
            os.makedirs("agg{0}D".format(i))
        file = open("agg{0}D\\agg{0}d.go".format(i), 'w')
        file.write(make(i))
        entries += [i]

    file = open("aggregation.go", 'w')
    file.write(make_aggregation(entries))


if __name__ == '__main__':
    main()
