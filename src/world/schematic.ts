export default class Schematic {
	public blocks: Array<number>;
	public data: Array<number>;

	constructor(public xSize: number, public ySize: number, public zSize: number) {
		this.blocks = new Array(xSize * ySize * zSize).fill(0);
		this.data = new Array(this.blocks.length).fill(0);
	}

	has(x: number, y: number, z: number): boolean {
		return this.blocks[this.getIndex(x, y, z)] !== 0;
	}

	get(x: number, y: number, z: number): [number, number] {
		const index = this.getIndex(x, y, z);
		return [
			this.blocks[index],
			this.data[index],
		];
	}

	set(x: number, y: number, z: number, block: number, data: number = 0) {
		const index = this.getIndex(x, y, z);
		this.blocks[index] = block;
		this.data[index] = data;
	}

	getIndex(x: number, y: number, z: number): number {
		return (y * this.zSize + z) * this.xSize + x;
	}
}
