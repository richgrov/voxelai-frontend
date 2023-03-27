export default class Schematic {
	private data: Array<boolean>;

	constructor(public xSize: number, public ySize: number, public zSize: number) {
		this.data = new Array(xSize * ySize * zSize);
	}

	has(x: number, y: number, z: number): boolean {
		return this.data[this.getIndex(x, y, z)];
	}

	set(x: number, y: number, z: number, value: boolean) {
		this.data[this.getIndex(x, y, z)] = value;
	}

	getIndex(x: number, y: number, z: number): number {
		return (y * this.ySize + z) * this.xSize + x;
	}
}
