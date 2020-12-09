package com.example.gdsudao.adapter

import android.view.LayoutInflater
import android.view.ViewGroup
import android.widget.ImageView
import android.widget.TextView
import androidx.recyclerview.widget.RecyclerView
import com.example.gdsudao.R
import com.example.gdsudao.model.Area

class AreaViewHolder(inflater: LayoutInflater, parent: ViewGroup): RecyclerView.ViewHolder(inflater.inflate(R.layout.adapter_area, parent, false)){
    private var tvNome: TextView? = null
    private var tvDataPlantio: TextView? = null
    private var tvDataPrevisao: TextView? = null
    private var tvNumeroCortes: TextView? = null
    private var tvProgresso: TextView? = null
    private var ivIcone: ImageView? = null

    init {
        tvNome = itemView.findViewById(R.id.tvNomeArea)
        tvDataPlantio = itemView.findViewById(R.id.tvDataPlantio)
        tvDataPrevisao = itemView.findViewById(R.id.tvDataPrevisao)
        tvNumeroCortes = itemView.findViewById(R.id.tvNumeroCortes)
        tvProgresso = itemView.findViewById(R.id.tvProgresso)
        ivIcone = itemView.findViewById(R.id.ivIconeSudao)
    }

    fun bind(area: Area){
        tvNome?.text = area.nome
        tvDataPlantio?.text = "Data Plantio: " + area.dataCorte
        //tvDataPrevisao?.text = "Próximo Corte: " + area.dataEstimada
        //tvNumeroCortes?.text = "Número de cortes: " + area.numeroCortes
        //tvProgresso?.text = "Data Plantio: " + area.progresso + "%"

        /*
        Implementar o esquema da figura
         */
    }
}

