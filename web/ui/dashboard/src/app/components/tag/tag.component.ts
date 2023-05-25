import { Component, Input, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { STATUS_COLOR } from 'src/app/models/global.model';

@Component({
	selector: 'convoy-tag, [convoy-tag]',
	standalone: true,
	imports: [CommonModule],
	template: `
		<ng-content></ng-content>
	`,
	styleUrls: ['./tag.component.scss'],
	host: { class: 'rounded-22px w-fit text-center text-12 justify-between gap-x-4px disabled:opacity-50', '[class]': 'classes' }
})
export class TagComponent implements OnInit {
	@Input('type') type: STATUS_COLOR = 'gray';
	@Input('fill') fill: 'outline' | 'soft' | 'solid' | 'soft-outline' = 'soft';
	@Input('color') color: 'primary' | 'danger' | 'warning' | 'gray' | 'success' = 'primary';
	@Input('size') size: 'sm' | 'md' | 'lg' = 'md';
	@Input('className') class!: string;

	fontSizes = { sm: 'text-10', md: `text-12`, lg: `text-14` };
	spacing = { sm: `py-2px px-6px`, md: `py-4px px-12px`, lg: `py-8px px-14px w-full` };

	constructor() {}

	ngOnInit(): void {}

	get classes() {
		const buttonTypes = {
			solid: `bg-${this.color}-400 text-white-100`,
			outline: `border border-${this.color}-400 text-${this.color}-400`,
			soft: `rounded-22px bg-${this.color}-50 text-${this.color}-400`,
			'soft-outline': `rounded-22px border-${this.color}-400 bg-${this.color}-50 text-${this.color}-400`
		};
		return `${this.fontSizes[this.size]} ${this.spacing[this.size]} ${buttonTypes[this.fill]}`;
	}
}
