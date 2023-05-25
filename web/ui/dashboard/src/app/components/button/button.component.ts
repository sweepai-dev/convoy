import { ChangeDetectionStrategy, Component, Input, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
	selector: 'convoy-button, [convoy-button]',
	standalone: true,
	host: { class: 'flex items-center justify-center disabled:opacity-50 cursor-pointer rounded-8px', '[class]': 'classes' },
	imports: [CommonModule],
	templateUrl: './button.component.html',
	changeDetection: ChangeDetectionStrategy.OnPush
})
export class ButtonComponent implements OnInit {
	@Input('buttonText') buttonText!: string;
	@Input('fill') fill: 'solid' | 'outline' | 'link' | 'ghost' | 'soft' = 'solid';
	@Input('size') size: 'xs' | 'sm' | 'md' | 'lg' = 'md';
	@Input('color') color: 'primary' | 'success' | 'warning' | 'danger' | 'grey' | 'transparent' | 'error' | 'gray' = 'primary';
	@Input('texture') texture: 'default' | 'light' = 'default';
	@Input('index') tabIndex = 0;
	@Input('spacing') spacing: 'true' | 'false' = 'true';
	fontSizes = { xs: 'text-12', sm: 'text-12', md: `text-14`, lg: `text-14` };
	buttonSpacing = { xs: 'py-4px px-8px', sm: `p-10px`, md: `py-12px px-16px`, lg: `py-18px px-20px w-full` };

	constructor() {}

	ngOnInit(): void {}

	get classes(): string {
		const colorLevel = this.texture == 'default' ? '400' : '25';
		const buttonTypes = {
			solid: `bg-${this.color}-${colorLevel} text-white-100 border-none`,
			soft: `bg-${this.color}-25 text-${this.color}-400 border-none border border-${this.color}-25`,
			outline: `border border-${this.color}-${colorLevel} text-${this.color}-400`,
			ghost: `border-none text-${this.color}-400`,
			link: `border-none text-${this.color}-400 underline decoration-${this.color}-400`
		};
		return `${this.fontSizes[this.size]} ${this.spacing == 'true' ? this.buttonSpacing[this.size] : ''} ${buttonTypes[this.fill]} flex items-center justify-center disabled:opacity-50`;
	}
}
