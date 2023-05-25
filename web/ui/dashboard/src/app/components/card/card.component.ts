import { Component, Input, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
	selector: '[convoy-card]',
	standalone: true,
	imports: [CommonModule],
	host: { class: 'rounded-8px transition-all duration-300 bg-white-100 border border-primary-25', '[class]': 'classes' },
	template: `
		<ng-content></ng-content>
	`
})
export class CardComponent implements OnInit {
	@Input('hover') hover: 'true' | 'false' = 'false';

	constructor() {}

	ngOnInit(): void {}

	get classes(): string {
		return `${this.hover === 'true' ? 'focus:shadow-sm hover:shadow-sm focus-visible:shadow-sm focus:border-grey-20 outline-none' : ''} block`;
	}
}
